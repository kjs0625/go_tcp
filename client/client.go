package client

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"reflect"
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

type Client struct {
	address string
	conn    net.Conn
}

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

func NewClient(address string) (*Client, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect server: %v", err)
	}

	return &Client{address: address, conn: conn}, nil
}

func (c *Client) SendPacket(packetType uint16, data []byte) error {

	header := &PacketHeader{}
	header.UiPacketSize = uint32(len(data))
	header.UsPacketType = packetType

	headerBuf := new(bytes.Buffer)
	if err := binary.Write(headerBuf, binary.LittleEndian, header); err != nil {
		return fmt.Errorf("Failted to serialize header: %v", err)
	}

	packetBuf := append(headerBuf.Bytes(), data...)
	size := uint32(binary.Size(header)) + header.UiPacketSize

	sent := 0
	for sent < int(size) {
		n, err := c.conn.Write(packetBuf[sent:size])
		if err != nil {
			if err == io.EOF {
				return fmt.Errorf("connection closed while sending")
			}
			return fmt.Errorf("send error: %v", err)
		}
		sent += n
	}

	return nil
}

func (c *Client) Close() {
	c.conn.Close()
}
