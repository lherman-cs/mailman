package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"

	"github.com/urfave/cli"
)

func receiverHandler(c *cli.Context) error {
	// Open receiver tcp port
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", Port))
	if err != nil {
		return err
	}
	defer listener.Close()

	ctx, cancel := context.WithCancel(context.Background())
	broadcastIP(ctx)

	conn, err := listener.Accept()
	cancel()
	if err != nil {
		return err
	}
	defer conn.Close()

	opt := Option{
		Interval:       256,
		LoadingMessage: "receiving",
		FinishMessage:  "received!",
	}
	_, err = io.CopyBuffer(os.Stdout, NewReader(conn, opt), nil)
	return err
}
