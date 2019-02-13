package main

import (
	"log"
	"os"
)

func main() {
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		log.Println("Sender mode")
		RunSender()
	} else {
		log.Println("Receiver mode")
		RunReceiver()
	}
}
