package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"syscall"
	"time"

	"github.com/grandcat/zeroconf"
	"github.com/manifoldco/promptui"
)

func RunSender() {
	entries := make([]*zeroconf.ServiceEntry, 0)
	hostnames := make([]string, 0)

	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		log.Fatalln("Failed to initialize resolver:", err.Error())
	}
	// Make a channel for results and start listening
	entriesCh := make(chan *zeroconf.ServiceEntry)

	go func() {
		for entry := range entriesCh {
			entries = append(entries, entry)
			hostnames = append(hostnames, entry.HostName)
		}
	}()

	// Start the lookup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	err = resolver.Browse(ctx, ServiceName, "local.", entriesCh)
	if err != nil {
		log.Fatalln("Failed to browse:", err.Error())
	}

	<-ctx.Done()
	cancel()

	prompt := promptui.Select{
		Label: "Send to",
		Items: hostnames,
	}

	tty, err := syscall.Open("/dev/tty", 0, 0)
	if err != nil {
		fmt.Println("Failed in creating tty")
		return
	}

	tmpFd, err := syscall.Dup(0)
	if err != nil {
		fmt.Println("Failed in creating a new temporary stdin fd")
		return
	}
	defer syscall.Close(tmpFd)

	// Swap stdin with new tty
	if err = syscall.Dup2(tty, 0); err != nil {
		fmt.Println("Failed in redirecting current stdin to new tty")
		return
	}

	idx, _, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	// Swapping back old stdin
	if err = syscall.Dup2(tmpFd, 0); err != nil {
		fmt.Println("Failed in swapping back old stdin")
		return
	}

	entry := entries[idx]
	if err = SendTo(fmt.Sprintf("%s:%d", entry.HostName, entry.Port)); err != nil {
		fmt.Println("Failed in sending the data to the receiver")
	}
}

func SendTo(address string) error {
	fmt.Println("Sending to", address)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}
	defer conn.Close()

	var ans rune
	fmt.Fscanf(conn, "%c", &ans)

	_, err = io.Copy(conn, os.Stdin)
	return err
}
