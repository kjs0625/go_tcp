package main

import (
	"fmt"
	"log"

	"tcp/client"
	"tcp/protocol"
)

func main() {
	address := "127.0.0.1:8080"
	c, err := client.NewClient(address)
	if err != nil {
		log.Fatalf("Failed to connect server: %v", err)
	}
	defer c.Close()
	fmt.Println("Success to connect server: ", address)

	payload := make([]byte, 1024)
	copy(payload, "Hello")
	err = c.SendPacket(protocol.SERVICE_1, payload)
	if err != nil {
		log.Fatalf("Failed to Send Packet : %v", err)
	}
}
