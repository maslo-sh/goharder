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
	"time"
)

const MaxBufferSize = 16384

type ProxyConfiguration struct {
	Name             string
	ListeningAddress model.Address
	Target           model.DataSource
	Sessions         map[string]*Session
	NumberOfSessions int
	Done             chan interface{}
}

type Session struct {
	Ctx        context.Context
	ClientConn net.Conn
	TargetConn net.Conn
}

func NewProxy(dto model.ProxyDto, ds model.DataSource) *ProxyConfiguration {
	return &ProxyConfiguration{
		Name:             dto.Name,
		ListeningAddress: dto.Address,
		Target:           ds,
		Sessions:         make(map[string]*Session),
		NumberOfSessions: 0,
	}
}

func (p *ProxyConfiguration) newSession(clientConn, targetConn net.Conn) string {
	sessionId := uuid.New()
	s := &Session{context.Background(), clientConn, targetConn}
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

func (p *ProxyConfiguration) Stop() {
	p.Done <- struct{}{}
}

func (p *ProxyConfiguration) CloseSessions() {
	for _, v := range p.Sessions {
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

	//defer listener.Close()

	log.Printf("Proxy listening on %s, forwarding to %s:%s...\n", p.ListeningAddress.CreateHostString(), p.Target.Hostname, p.Target.Port)

	for {
		select {
		case <-p.Done:
			log.Printf("Closed listener for proxy")
			return
		default:
			clientConn, err := listener.Accept()
			if err != nil {
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
}

func (p *ProxyConfiguration) handleConnection(session *Session) {
	clientBuffer := make([]byte, MaxBufferSize)
	targetBuffer := make([]byte, MaxBufferSize)

	go handleSession(session.Ctx, session.ClientConn, session.TargetConn, clientBuffer)
	go handleSession(session.Ctx, session.TargetConn, session.ClientConn, targetBuffer)

	for {
		select {
		case <-session.Ctx.Done():
			log.Printf("Proxy has been stopped")
			session.Close()
		default:
			// Perform task operation
			log.Println("Performing task...")
			time.Sleep(10 * time.Second)
		}
	}
}

func handleSession(ctx context.Context, srcConn net.Conn, dstConn net.Conn, buff []byte) {
	var packetsProcessed int
	//ctx, cancel := context.WithTimeout(ctx, time.Second*10)

	for {
		read, err := srcConn.Read(buff)
		if err != nil {
			if err != io.EOF {
				log.Printf("Error reading data from source: %v\n", err)
			}
			break
		}

		packetsProcessed = packetsProcessed + 1

		written, err := dstConn.Write(buff[:read])

		if err != nil {
			if err != io.EOF {
				log.Printf("Error writing data from target to destination: %v\n", err)
			}
			break
		}

		if packetsProcessed == 1 {
			log.Printf("It's a startup message")
			messageLength := binary.BigEndian.Uint32(buff[:4])
			log.Printf("STARTUP MESSAGE LENGTH: %d", messageLength)
			continue
		}

		next := true
		offset := 0
		for next {
			header := CreateHeaderFromBytes(buff[offset : offset+5])
			log.Printf("Packet type: %s; packet length: %d", GetPacketType(header.PacketType), header.PacketLength)
			offset = offset + header.PacketLength + 1
			next = !(offset >= read)
		}

		log.Printf("Copied %d bytes from listener to target\n", written)
	}
}
