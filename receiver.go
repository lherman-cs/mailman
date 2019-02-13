package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/hashicorp/mdns"
)

func RunReceiver() {
	var ans rune

	list, err := net.Listen("tcp", fmt.Sprintf(":%d", Port))
	if err != nil {
		log.Fatal(err)
	}

	// Setup our service export
	host, _ := os.Hostname()
	info := []string{"Mailman"}
	service, _ := mdns.NewMDNSService(host, ServiceName, "", "", Port, nil, info)

	// Create the mDNS server, defer shutdown
	server, _ := mdns.NewServer(&mdns.Config{Zone: service})
	defer server.Shutdown()

	for {
		conn, err := list.Accept()
		if err != nil {
			continue
		}

		log.Printf("Do you want to receive from %s? [y/n] ", conn.LocalAddr())
		fmt.Scanf("%c", &ans)

		// send back answer
		fmt.Fprintf(conn, "%c", ans)

		if ans == 'y' {
			io.Copy(os.Stdout, conn)
			break
		}

		conn.Close()
	}
}
