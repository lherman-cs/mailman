package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = name
	app.Usage = "Peer-to-peer Data Sharing"
	app.Version = version

	app.Commands = []cli.Command{
		{
			Name:    "send",
			Aliases: []string{"s"},
			Usage:   "send data from stdin to the other peer",
			Action:  SenderHandler,
		},
		{
			Name:    "receive",
			Aliases: []string{"r", "recv"},
			Usage:   "receive data and output to stdout",
			Action:  ReceiverHandler,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
