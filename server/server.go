package server

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"reflect"
	"sync"
	"time"
)

const (
	SERVICE_1 = (10001)
	SERVICE_2 = (10002)
	SERVICE_3 = (10003)
	SERVICE_4 = (10004)
)

type PacketHeader struct {
	UiPacketSize uint32
	UsPacketType uint16
}

type Server struct {
	address     string
	listener    net.Listener
	mu          sync.Mutex
	clients     []net.Conn
	readTimeout time.Duration
}

// util
func SizeOf(value any) (size int) {
	t := reflect.TypeOf(value)
	v := reflect.ValueOf(value)

	switch t.Kind() {
	case reflect.Array:
		elem := t.Elem()
		size = int(elem.Size()) * v.Len()
	case reflect.Struct:
		sum := 0
		for i, n := 0, v.NumField(); i < n; i++ {
			s := SizeOf(v.Field(i).Interface())
			if s < 0 {
				break
			}
			sum += s
		}
		size = sum
	case reflect.Int:
		size = 4
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
		size = int(t.Size())
	case reflect.Slice:
		size = v.Cap()
	case reflect.Bool:
		size = 1
	default:
	}

	return size
}

func NewServer(address string, readTimeout time.Duration) *Server {
	return &Server{
		address:     address,
		readTimeout: readTimeout,
		clients:     make([]net.Conn, 0),
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
		s.clients = append(s.clients, conn)
		s.mu.Unlock()

		go s.Handler(conn)
	}
}

func (s *Server) Handler(conn net.Conn) {
	defer conn.Close()

	for {
		if err := conn.SetReadDeadline(time.Now().Add(s.readTimeout)); err != nil {
			fmt.Println("Error setting read deadline:", err)
			return
		}

		header := &PacketHeader{}
		headerBuffer := make([]byte, SizeOf(PacketHeader{}))

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

		s.HandleServicePacket(header, payloadBuffer)
	}
}

func (s *Server) HandleServicePacket(header *PacketHeader, payload []byte) {
	switch header.UsPacketType {
	case SERVICE_1:
		s.HandleService1(payload)
	case SERVICE_2:
		s.HandleService2(payload)
	case SERVICE_3:
		s.HandleService3(payload)
	}
}

func (s *Server) HandleService1(payLoad []byte) {
	fmt.Println("HandleService1")
}

func (s *Server) HandleService2(payLoad []byte) {
	fmt.Println("HandleService2")
}

func (s *Server) HandleService3(payLoad []byte) {
	fmt.Println("HandleService3")
}

func (s *Server) Stop() {
	fmt.Println("Shutting down server...")
	s.mu.Lock()
	for _, conn := range s.clients {
		conn.Close()
	}
	s.listener.Close()
	s.mu.Unlock()
	fmt.Println("Server shut down")
}
