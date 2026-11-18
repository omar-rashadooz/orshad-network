package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

var dnsRegistry = map[string]string{
	"linux.osh":  "127.0.0.1:8585",
	"games.osh": "127.0.0.1:8585",
	"news.osh":  "127.0.0.1:8585",
	"orshad.osh":  "127.0.0.1:8585",
}

func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:5353")
	if err != nil {
		return
	}
	defer listener.Close()
	fmt.Println("=== ORSHAD DNS SERVER OPERATIONAL ON PORT 5353 ===")

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleDNSQuery(conn)
	}
}

func handleDNSQuery(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	domain, err := reader.ReadString('\n')
	if err != nil {
		return
	}

	domain = strings.TrimSpace(domain)
	address, exists := dnsRegistry[domain]

	if exists {
		conn.Write([]byte(address + "\n"))
	} else {
		conn.Write([]byte("NXDOMAIN\n"))
	}
}
