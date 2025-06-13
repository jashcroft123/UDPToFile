package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", "224.0.0.1:9999")
	if err != nil {
		fmt.Println("ResolveUDPAddr failed:", err)
		os.Exit(1)
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		fmt.Println("DialUDP failed:", err)
		os.Exit(1)
	}
	defer conn.Close()

	msg := []byte("Hello from multicast sender!")

	for {
		_, err := conn.Write(msg)
		if err != nil {
			fmt.Println("Write failed:", err)
		} else {
			fmt.Println("Sent:", string(msg))
		}
		time.Sleep(2 * time.Second)
	}
}
