package main

import (
	"fmt"
	"net"
	"time"
)

func dial(network string, address string, retry int) (net.Conn, error) {
	for i := 0; i < retry; i++ {
		if conn, err := net.Dial(network, address); err == nil {
			return conn, nil
		}
		// TODO! Probably more flexible
		time.Sleep(time.Second)
	}
	return nil, fmt.Errorf("Failed after %d retries", retry)
}
