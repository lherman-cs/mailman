package main

import (
	"log"

	"github.com/phayes/freeport"
)

const (
	ServiceName = "_mailman._tcp"
)

var (
	Port int
)

func init() {
	var err error
	Port, err = freeport.GetFreePort()
	if err != nil {
		log.Fatal(err)
	}
}
