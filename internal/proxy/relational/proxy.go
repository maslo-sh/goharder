package relational

import (
	"context"
	"encoding/binary"
	"fmt"
	"github.com/bluele/gcache"
	"io"
	"log"
	"math"
	"net"
	utils "proxy-engineering-thesis"
	"proxy-engineering-thesis/internal/aws"
	"proxy-engineering-thesis/internal/proxy/detection"
	"proxy-engineering-thesis/model"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	MaxBufferSize = 16400
	DetectionMode = iota
	PreventionMode
	FullProtectionMode
)

type ProxyConfiguration struct {
	Name             string
	ListeningAddress model.Address
	Listener         net.Listener
	Target           model.DataSource
	Sessions         map[string]*Session
	NumberOfSessions int
	Mode             int
	Done             chan interface{}
}

type Session struct {
	Ctx                       context.Context
	ClientConn                net.Conn
	TargetConn                net.Conn
	detector                  detection.Detector
	BlackListCache            gcache.Cache
	ClosingTriggered          bool
	ClientActivityInterrupted bool
	mode                      int
	cwClient                  *aws.CloudWatchConfiguration
}

func NewProxy(dto model.ProxyDto, ds model.DataSource, proxyMode string) *ProxyConfiguration {
	return &ProxyConfiguration{
		Name:             dto.Name,
		ListeningAddress: dto.Address,
		Target:           ds,
		Sessions:         make(map[string]*Session),
		NumberOfSessions: 0,
		Mode:             GetProxyMode(proxyMode),
		Done:             make(chan interface{}),
	}
}

func (p *ProxyConfiguration) newSession(clientConn, targetConn net.Conn) string {
	var cwClient *aws.CloudWatchConfiguration
	conf, err := utils.ReadPropertiesBasedConfig("resources/cw.properties")
	if err != nil {
		fmt.Printf("failed to retrieve CloudWatch config from file: %v\n", err)
	}

	if conf["logGroupName"] == "" || conf["logStreamName"] == "" {
		fmt.Printf("didn't found necessary CloudWatch configuration in resources")
		cwClient = nil
	} else {
		cwClient = aws.NewCloudWatchConfiguration(conf["logGroupName"], conf["logStreamName"])
	}
	sessionId := uuid.New()
	s := &Session{
		context.Background(),
		clientConn,
		targetConn,
		detection.SqlDetector{},
		gcache.New(100).LFU().LoaderFunc(func(key interface{}) (interface{}, error) {
			return 0, nil
		}).EvictedFunc(func(key, val interface{}) {
			strKey, ok := key.(string)
			if ok {
				log.Printf("evicted entry from cache - %s can connect again without delay", strKey)
			}

		}).Build(),
		false,
		false,
		p.Mode,
		cwClient,
	}
	go s.cwClient.InitLogStore()
	p.Sessions[sessionId.String()] = s
	p.NumberOfSessions = p.NumberOfSessions + 1
	log.Printf("created session: %s", sessionId)
	return sessionId.String()
}

func (s *Session) Close() error {

	err := s.ClientConn.Close()
	if err != nil {
		log.Printf("failed to close client connection: %v\n", err)
		return err
	}
	err = s.TargetConn.Close()
	if err != nil {
		log.Printf("failed to close target connection: %v\n", err)
		return err
	}
	s.BlackListCache.Purge()
	return nil
}

func (p *ProxyConfiguration) CloseSessions() {
	for _, v := range p.Sessions {
		v.ClosingTriggered = true
		v.Close()
	}
}

func (p *ProxyConfiguration) Start() {
	log.Printf("started new '%s' proxy instance", p.Name)
	listener, err := net.Listen("tcp", p.ListeningAddress.CreateHostString())
	if err != nil {
		log.Printf("error when setting up listener: %v\n", err)
		return
	}
	p.Listener = listener

	defer listener.Close()

	log.Printf("proxy listening on %s, forwarding to %s:%s...\n", p.ListeningAddress.CreateHostString(), p.Target.Hostname, p.Target.Port)

	for {
		clientConn, err := listener.Accept()
		if err != nil {
			if strings.Contains(err.Error(), "use of closed network connection") {
				log.Printf("closed proxy listener")
				p.Done <- struct{}{}
				return
			}
			log.Printf("error accepting connection: %v\n", err)
			continue
		}

		log.Printf("received new connection: %s", clientConn.RemoteAddr().String())

		targetConn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", p.Target.Hostname, p.Target.Port))
		if err != nil {
			log.Printf("target is unreachable: %v\n", err)
			clientConn.Close()
			continue
		}

		log.Printf("set up connection with database")

		sessionId := p.newSession(clientConn, targetConn)
		session := p.Sessions[sessionId]

		go p.handleConnection(session)
	}
}

func (p *ProxyConfiguration) handleConnection(session *Session) {
	var packetsProcessed int
	clientBuffer := make([]byte, MaxBufferSize)
	targetBuffer := make([]byte, MaxBufferSize)

	for {
		if session.ClosingTriggered {
			p.Done <- struct{}{}
		}

		select {
		case <-p.Done:
			return
		default:

			packetsProcessed = packetsProcessed + 1
			session.ClientActivityInterrupted = false

			// Handle potential blacklisting
			source := session.ClientConn.RemoteAddr().String()
			cacheValue, _ := session.BlackListCache.Get(source)
			exp, ok := cacheValue.(int)
			if !ok {
				log.Printf("failed to convert cache value under '%s' key - value is not an integer", source)
			} else {
				timeOfSleeping := math.Pow(2, float64(exp))
				time.Sleep(time.Duration(timeOfSleeping) * time.Second)
			}

			// Perform task operation
			session.handleInboundTraffic(clientBuffer, packetsProcessed)
			if !session.ClientActivityInterrupted {
				session.handleOutboundTraffic(targetBuffer, packetsProcessed)
			}
		}
	}
}

func (s *Session) UpdateBlackList(key string, value interface{}) {
	intValue, ok := value.(int)
	if !ok {
		log.Printf("failed to load cache value under '%s' key - value is not an integer", key)
		return
	}
	exp := math.Min(float64(intValue)+1, 6)
	newBlockTime := time.Second * time.Duration(math.Pow(2, exp+1))
	s.BlackListCache.SetWithExpire(key, exp, newBlockTime)
	log.Printf("next request from source %s will be blocked for %d seconds", key, newBlockTime/1e9)
}

func (s *Session) handleInboundTraffic(buff []byte, packetsProcessed int) {
	read, err := s.ClientConn.Read(buff)
	if err != nil {
		if strings.Contains(err.Error(), "use of closed network connection") {
			log.Printf("Connection closed")
			s.ClosingTriggered = true
			return
		}
		if err != io.EOF {
			log.Printf("Error reading data from source: %v\n", err)
		}
		return
	}

	next := true
	maliciousDetected := false
	offset := 0
	for next {
		header := CreateHeaderFromBytes(buff[offset : offset+5])
		if header.PacketType == 0x51 || header.PacketType == 0x50 {
			status := s.detector.DetectMaliciousContent(buff[offset+5 : offset+header.PacketLength])
			if status == detection.MALICIOUS {
				maliciousDetected = true
			}
		}
		log.Printf("Packet type: %s; packet length: %d", GetPacketType(header.PacketType), header.PacketLength)
		offset = offset + header.PacketLength + 1
		next = !(offset >= read)
	}

	var buffToWrite []byte
	if maliciousDetected {
		source := s.ClientConn.RemoteAddr().String()
		numberOfAttempts, err := s.BlackListCache.Get(source)
		if err != nil {
			log.Printf("failed to retrieve value under key %s from cache: %v\n", source, err)
		} else {
			s.UpdateBlackList(source, numberOfAttempts)
		}

		if s.mode == DetectionMode || s.mode == FullProtectionMode {
			message := fmt.Sprintf("Malicious query; IP - %s", s.ClientConn.RemoteAddr())
			go s.cwClient.SendLog(message)
		}

		if s.mode == PreventionMode || s.mode == FullProtectionMode {
			buffToWrite = GetMaliciousActivityDetectedError()
			_, err := s.ClientConn.Write(buffToWrite)
			if err != nil {
				if err != io.EOF {
					log.Printf("Error writing data to proxy message buffer: %v\n", err)
				}
			}
			s.ClientActivityInterrupted = true
			return
		}
	}

	buffToWrite = buff[:read]
	written, err := s.TargetConn.Write(buffToWrite)

	if err != nil {
		if strings.Contains(err.Error(), "use of closed network connection") {
			log.Printf("Connection closed")
			s.ClosingTriggered = true
			return
		}
		if err != io.EOF {
			log.Printf("Error writing data from target to destination: %v\n", err)
		}
		return
	}

	if packetsProcessed == 1 {
		log.Printf("It's a startup message")
		messageLength := binary.BigEndian.Uint32(buff[:4])
		log.Printf("STARTUP MESSAGE LENGTH: %d", messageLength)
		return
	}

	log.Printf("Copied %d bytes from listener to target\n", written)
}

func (s *Session) handleOutboundTraffic(buff []byte, packetsProcessed int) {
	//when sending message from proxy itself, program will lock on that reading command
	read, err := s.TargetConn.Read(buff)
	if err != nil {
		if strings.Contains(err.Error(), "use of closed network connection") {
			log.Printf("Connection closed")
			return
		}
		if err != io.EOF {
			log.Printf("Error reading data from source: %v\n", err)
		}
		return
	}

	written, err := s.ClientConn.Write(buff[:read])

	if err != nil {
		if strings.Contains(err.Error(), "use of closed network connection") {
			log.Printf("Connection closed")
			return
		}
		if err != io.EOF {
			log.Printf("Error writing data from target to destination: %v\n", err)
		}
		return
	}

	if packetsProcessed == 1 {
		log.Printf("It's a startup message")
		messageLength := binary.BigEndian.Uint32(buff[:4])
		log.Printf("STARTUP MESSAGE LENGTH: %d", messageLength)
		return
	}

	readPacketBytes(buff, read)

	log.Printf("Copied %d bytes from listener to target\n", written)
}

func readPacketBytes(buff []byte, readBytes int) {
	next := true
	offset := 0
	for next {
		header := CreateHeaderFromBytes(buff[offset : offset+5])
		log.Printf("Packet type: %s; packet length: %d", GetPacketType(header.PacketType), header.PacketLength)
		offset = offset + header.PacketLength + 1
		next = !(offset >= readBytes)
	}
}

func GetProxyMode(mode string) int {
	lowerCasedMode := strings.ToLower(mode)
	if lowerCasedMode == "prevention" {
		return PreventionMode
	}
	if lowerCasedMode == "detection" {
		return DetectionMode
	}

	return FullProtectionMode

}
