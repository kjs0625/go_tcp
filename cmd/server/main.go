package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"tcp/server"
)

func main() {
	address := ":8080"
	readTimeout := 30 * time.Second

	srv := server.NewServer(address, readTimeout)

	go func() {
		fmt.Println("Starting TCP Server...")
		if err := srv.Start(); err != nil {
			fmt.Printf("Error starting server: %v\n", err)
		}
	}()

	// Graceful Shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	sig := <-sigChan
	switch sig {
	case os.Interrupt:
		fmt.Println("Exit with os.Interrupt")
	case syscall.SIGTERM:
		fmt.Println("Exit with syscall.SIGTERM")
	}

	fmt.Println("Shutdown signal received...")
	srv.Stop()

	fmt.Println("Server exited...")
}
