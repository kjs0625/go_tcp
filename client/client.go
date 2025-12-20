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

func (c *Client) NewClient(address string) (*Client, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect server: %v", err)
	}

	return &Client{address: address, conn: conn}, nil
}

func (c *Client) Handler(conn net.Conn) {
	defer conn.Close()

	for {
		header := &protocol.PacketHeader{}
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

		if header.UiPacketSize > protocol.MAX_PACKET_SIZE {
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

		c.HandleServicePacket(header.UsPacketType, payloadBuffer)
	}
}

func (c *Client) HandleServicePacket(packetType uint16, payload []byte) {
	switch packetType {
	case protocol.SERVICE_1_RESP:
		c.HandleService1Response(payload)
	case protocol.SERVICE_2_RESP:
		c.HandleService2Response(payload)
	case protocol.SERVICE_3_RESP:
		c.HandleService3Response(payload)
	default:
		fmt.Println("Unknown packet type: ", packetType)
	}
}

func (c *Client) HandleService1Response(payload []byte) {
	fmt.Println("HandleService1 Response")
	fmt.Println(string(payload))
}

func (c *Client) HandleService2Response(payload []byte) {
	fmt.Println("HandleService2 Response")
	fmt.Println(string(payload))
}

func (c *Client) HandleService3Response(payload []byte) {
	fmt.Println("HandleService1")
	fmt.Println(string(payload))
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
