package main

import (
	"context"
	"fmt"
	"net"
	"time"
)

type ReceiverInfo struct {
	IP       string
	Hostname string
}

func getReceiverInfo() (*ReceiverInfo, error) {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", Port))
	if err != nil {
		return nil, err
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	buf := make([]byte, 1)
	_, receiver, err := conn.ReadFromUDP(buf)
	ip := receiver.IP.String()
	hosts, err := net.LookupAddr(ip)
	hostname := "-"
	if err == nil && len(hosts) > 0 {
		hostname = hosts[len(hosts)-1]
	}

	info := ReceiverInfo{
		IP:       receiver.IP.String(),
		Hostname: hostname,
	}
	return &info, err
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
