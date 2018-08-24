package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "Share"
	app.Usage = "Peer-to-peer share data"
	app.Version = "0.1.2"

	app.Commands = []cli.Command{
		{
			Name:    "send",
			Aliases: []string{"s"},
			Usage:   "get from stdin out to `share`",
			Action:  senderHandler,
		},
		{
			Name:    "receive",
			Aliases: []string{"r", "recv"},
			Usage:   "get from `share` out to stdout",
			Action:  receiverHandler,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
