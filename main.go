package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/urfave/cli"
)

func senderHandler(c *cli.Context) error {
	ip, err := getReceiverIP()
	if err != nil {
		return err
	}
	log.Println("Found", ip)

	conn, err := dial("tcp", fmt.Sprintf("%s:%d", ip, Port), 5)
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

func receiverHandler(c *cli.Context) error {
	// Open receiver tcp port
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", Port))
	if err != nil {
		return err
	}
	defer listener.Close()

	if err := broadcastIP(); err != nil {
		return err
	}

	conn, err := listener.Accept()
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

func main() {
	app := cli.NewApp()
	app.Name = "Share"
	app.Usage = "Peer-to-peer share data"
	app.Version = "0.1.2"

	app.Commands = []cli.Command{
		{
			Name:   "send",
			Usage:  "get from stdin out to `share`",
			Action: senderHandler,
		},
		{
			Name:   "receive",
			Usage:  "get from `share` out to stdout",
			Action: receiverHandler,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}