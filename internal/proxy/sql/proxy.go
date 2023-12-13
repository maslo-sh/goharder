package sql

import (
	"context"
	"encoding/binary"
	"fmt"
	"github.com/google/uuid"
	"io"
	"log"
	"net"
	"proxy-engineering-thesis/model"
	"strings"
	"time"
)

const (
	MaxBufferSize = 16400
	DetectionMode = iota
	PreventionMode
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
	ClosingTriggered          bool
	ClientActivityInterrupted bool
}

func NewProxy(dto model.ProxyDto, ds model.DataSource) *ProxyConfiguration {
	return &ProxyConfiguration{
		Name:             dto.Name,
		ListeningAddress: dto.Address,
		Target:           ds,
		Sessions:         make(map[string]*Session),
		NumberOfSessions: 0,
		Done:             make(chan interface{}),
	}
}

func (p *ProxyConfiguration) newSession(clientConn, targetConn net.Conn) string {
	sessionId := uuid.New()
	s := &Session{
		context.Background(),
		clientConn,
		targetConn,
		false,
		false}
	p.Sessions[sessionId.String()] = s
	p.NumberOfSessions = p.NumberOfSessions + 1
	log.Printf("Created session: %s", sessionId)
	return sessionId.String()
}

func (s *Session) Close() error {
	err := s.ClientConn.Close()
	if err != nil {
		log.Printf("Failed to close client connection: %v\n", err)
		return err
	}
	err = s.TargetConn.Close()
	if err != nil {
		log.Printf("Failed to close target connection: %v\n", err)
		return err
	}
	return nil
}

func (p *ProxyConfiguration) RemoveSession(sessionId string) {
	delete(p.Sessions, sessionId)
	p.NumberOfSessions = p.NumberOfSessions - 1
	log.Printf("Removed session: %s", sessionId)
}

func (p *ProxyConfiguration) CloseSessions() {
	for _, v := range p.Sessions {
		v.ClosingTriggered = true
		v.Close()
	}
}

func (p *ProxyConfiguration) Start() {
	log.Printf("Started new '%s' proxy instance", p.Name)
	listener, err := net.Listen("tcp", p.ListeningAddress.CreateHostString())
	if err != nil {
		log.Printf("Error: %v\n", err)
		return
	}
	p.Listener = listener

	defer listener.Close()

	log.Printf("Proxy listening on %s, forwarding to %s:%s...\n", p.ListeningAddress.CreateHostString(), p.Target.Hostname, p.Target.Port)

	for {
		clientConn, err := listener.Accept()
		if err != nil {
			if strings.Contains(err.Error(), "use of closed network connection") {
				log.Printf("Closed proxy listener")
				p.Done <- struct{}{}
				return
			}
			log.Printf("Error accepting connection: %v\n", err)
			continue
		}

		log.Printf("Received new connection: %s", clientConn.RemoteAddr().String())

		targetConn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", p.Target.Hostname, p.Target.Port))
		if err != nil {
			log.Printf("Target is unreachable: %v\n", err)
			clientConn.Close()
			continue
		}

		log.Printf("Set up connection with database")

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
		select {
		case <-p.Done:
			return
		default:
			packetsProcessed = packetsProcessed + 1
			session.ClientActivityInterrupted = false
			// Perform task operation
			session.handleInboundTraffic(clientBuffer, packetsProcessed)
			if !session.ClientActivityInterrupted {
				session.handleOutboundTraffic(targetBuffer, packetsProcessed)
			}
			time.Sleep(10 * time.Millisecond)
		}
	}
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
			query := string(buff[offset+5 : offset+header.PacketLength])
			log.Printf("PACKET WITH DQL: %s", query)
			if strings.Contains(query, "OR'1=1';") {
				maliciousDetected = true
			}
		}
		log.Printf("Packet type: %s; packet length: %d", GetPacketType(header.PacketType), header.PacketLength)
		offset = offset + header.PacketLength + 1
		next = !(offset >= read)
	}

	var buffToWrite []byte
	if maliciousDetected {
		buffToWrite = GetMaliciousActivityDetectedError()
		_, err := s.ClientConn.Write(buffToWrite)
		if err != nil {
			if err != io.EOF {
				log.Printf("Error writing data to proxy message buffer: %v\n", err)
			}
		}
		s.ClientActivityInterrupted = true
		return
	} else {
		buffToWrite = buff[:read]
	}

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
