package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

const (
	multicastAddr = "234.5.6.7:50007" // Multicast group and port
	logFilePath   = "multicast.log"   // Log file path
)

func init() {
	fmt.Printf("Started listening on", multicastAddr)
}

func main() {
	// Set up logging to file
	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer file.Close()
	log.SetOutput(file)

	log.Println("Starting UDP multicast listener...")

	// Resolve address
	addr, err := net.ResolveUDPAddr("udp", multicastAddr)
	if err != nil {
		log.Fatalf("Failed to resolve UDP address: %v", err)
	}

	// Set up listener
	conn, err := net.ListenMulticastUDP("udp", nil, addr)
	if err != nil {
		log.Fatalf("Failed to join multicast group: %v", err)
	}
	defer conn.Close()

	buffer := make([]byte, 1024)

	for {
		n, src, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Printf("Error reading from socket: %v", err)
			continue
		}
		message := string(buffer[:n])
		log.Printf("Received from %v: %s", src, message)
		fmt.Printf("Received from %v: %s\n", src, message)
	}
}
