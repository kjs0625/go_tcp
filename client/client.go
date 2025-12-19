package client

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"

	"tcp/protocol"
)

type Client struct {
	address string
	conn    net.Conn
}

func NewClient(address string) (*Client, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect server: %v", err)
	}

	return &Client{address: address, conn: conn}, nil
}

func (c *Client) SendPacket(packetType uint16, data []byte) error {
	header := &protocol.PacketHeader{}
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
