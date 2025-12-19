package server

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

const (
	SERVICE_1 = (10001)
	SERVICE_2 = (10002)
	SERVICE_3 = (10003)
	SERVICE_4 = (10004)

	MAX_PACKET_SIZE = 1024 * 8
)

type PacketHeader struct {
	UiPacketSize uint32
	UsPacketType uint16
}

type Server struct {
	address     string
	listener    net.Listener
	mu          sync.Mutex
	clients     map[net.Conn]bool
	readTimeout time.Duration
}

func NewServer(address string, readTimeout time.Duration) *Server {
	return &Server{
		address:     address,
		readTimeout: readTimeout,
		clients:     make(map[net.Conn]bool), // init map
	}
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		return fmt.Errorf("Failed to start server: %v", err)
	}
	s.listener = listener
	fmt.Println("Server started on", s.address)

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection", err)
			continue
		}
		fmt.Println("Client connected:", conn.RemoteAddr())

		s.mu.Lock()
		s.clients[conn] = true
		s.mu.Unlock()

		go s.Handler(conn)
	}
}

func (s *Server) Handler(conn net.Conn) {
	// clean up
	defer func() {
		conn.Close()
		s.mu.Lock()
		delete(s.clients, conn) // delete from map
		s.mu.Unlock()
		fmt.Println("Client disconnected cleanup: ", conn.RemoteAddr())
	}()

	for {
		if err := conn.SetReadDeadline(time.Now().Add(s.readTimeout)); err != nil {
			fmt.Println("Error setting read deadline:", err)
			return
		}

		header := &PacketHeader{}
		headerBuffer := make([]byte, binary.Size(header))

		_, err := io.ReadFull(conn, headerBuffer)
		if err != nil {
			if err == io.EOF {
				fmt.Println("Client disconnected:", conn.RemoteAddr())
				return
			} else {
				fmt.Println("Error reading size header:", err)
				return
			}
		}
		rdr := bytes.NewReader(headerBuffer)
		if err := binary.Read(rdr, binary.LittleEndian, header); err != nil {
			fmt.Println("Failed to Read headerBuffer")
			return
		}

		if header.UiPacketSize > MAX_PACKET_SIZE {
			fmt.Printf("Error: Packet size too large (%d)\n", header.UiPacketSize)
			return
		}

		payloadBuffer := make([]byte, header.UiPacketSize)
		_, err = io.ReadFull(conn, payloadBuffer)
		if err != nil {
			if err == io.EOF {
				fmt.Println("Client disconnected:", conn.RemoteAddr())
				return
			} else {
				fmt.Println("Error reading data:", err)
				return
			}
		}

		s.HandleServicePacket(header.UsPacketType, payloadBuffer)
	}
}

func (s *Server) HandleServicePacket(packetType uint16, payload []byte) {
	switch packetType {
	case SERVICE_1:
		s.HandleService1(payload)
	case SERVICE_2:
		s.HandleService2(payload)
	case SERVICE_3:
		s.HandleService3(payload)
	default:
		fmt.Println("Unknown packet type: ", packetType)
	}
}

func (s *Server) HandleService1(payLoad []byte) {
	fmt.Println("HandleService1")
	fmt.Println(string(payLoad))
}

func (s *Server) HandleService2(payLoad []byte) {
	fmt.Println("HandleService2")
	fmt.Println(string(payLoad))
}

func (s *Server) HandleService3(payLoad []byte) {
	fmt.Println("HandleService3")
	fmt.Println(string(payLoad))
}

func (s *Server) Stop() {
	fmt.Println("Shutting down server...")
	s.mu.Lock()
	defer s.mu.Unlock()

	// close listener first
	if s.listener != nil {
		s.listener.Close()
	}

	// close conn's
	for conn := range s.clients {
		conn.Close()
	}
	s.clients = make(map[net.Conn]bool) // init map
	fmt.Println("Server shut down")
}
