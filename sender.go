package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/urfave/cli"
)

func SenderHandler(c *cli.Context) error {
	info, err := getReceiverInfo()
	if err != nil {
		return err
	}
	log.Printf("Found: %s (%s)", info.Hostname, info.IP)

	conn, err := dial("tcp", fmt.Sprintf("%s:%d", info.IP, Port), 5)
	if err != nil {
		return err
	}
	defer conn.Close()

	opt := Option{
		Interval:       256,
		LoadingMessage: "sending",
		FinishMessage:  "sent!",
	}
	_, err = io.CopyBuffer(conn, NewReader(os.Stdin, opt), nil)
	return err
}
