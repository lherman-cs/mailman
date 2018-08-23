package main

import (
	"context"
	"fmt"
	"net"
	"time"
)

func getReceiverIP() (string, error) {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", Port))
	if err != nil {
		return "", err
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	buf := make([]byte, 1) // HACK! FIXME
	_, receiver, err := conn.ReadFromUDP(buf)
	return receiver.IP.String(), err
}

func broadcastIP(ctx context.Context) <-chan error {
	errc := make(chan error)
	go func() {
		defer close(errc)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				conn, err := net.Dial("udp", fmt.Sprintf("255.255.255.255:%d", Port))
				if err != nil {
					errc <- err
					return
				}
				defer conn.Close()

				fmt.Fprintf(conn, "")
				time.Sleep(time.Second)
			}
		}
	}()
	return errc
}
