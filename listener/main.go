package main

import (
	"fmt"
	"net"
	"os"

	"golang.org/x/net/ipv4"
)

func main() {
	group := net.IPv4(224, 0, 0, 1)
	port := 9999

	// Listen on all interfaces
	addr := &net.UDPAddr{
		IP:   net.IPv4zero, // 0.0.0.0
		Port: port,
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("ListenUDP failed:", err)
		os.Exit(1)
	}
	defer conn.Close()

	p := ipv4.NewPacketConn(conn)

	// Enable multicast loopback so we receive our own packets (for loopback testing)
	if err := p.SetMulticastLoopback(true); err != nil {
		fmt.Println("SetMulticastLoopback failed:", err)
	}

	// Join group on all interfaces
	ifaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("Interfaces error:", err)
		os.Exit(1)
	}

	for _, iface := range ifaces {
		// Skip interfaces that are down or don't support multicast
		if iface.Flags&net.FlagUp != 0 && iface.Flags&net.FlagMulticast != 0 {
			err := p.JoinGroup(&iface, &net.UDPAddr{IP: group})
			if err != nil {
				fmt.Printf("Failed to join group on %s: %v\n", iface.Name, err)
			} else {
				fmt.Println("Joined group on interface:", iface.Name)
			}
		}
	}

	buf := make([]byte, 1500)
	fmt.Println("Listening for multicast messages...")

	for {
		n, _, src, err := p.ReadFrom(buf)
		if err != nil {
			fmt.Println("ReadFrom failed:", err)
			continue
		}
		fmt.Printf("Received from %v: %s\n", src, string(buf[:n]))
	}
}
