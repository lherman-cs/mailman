package main

import (
	"fmt"
	"net"
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

func broadcastIP() error {
	conn, err := net.Dial("udp", fmt.Sprintf("255.255.255.255:%d", Port))
	if err != nil {
		return err
	}
	defer conn.Close()

	fmt.Fprintf(conn, "a") // HACK! FIXME
	return nil
}
