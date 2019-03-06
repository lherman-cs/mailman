package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/grandcat/zeroconf"
)

func RunReceiver() {
	var ans rune

	list, err := net.Listen("tcp", fmt.Sprintf(":%d", Port))
	if err != nil {
		log.Fatal(err)
	}
	defer list.Close()

	server, err := zeroconf.Register("Mailman", ServiceName, "local.", Port, nil, nil)
	if err != nil {
		panic(err)
	}
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
