package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"syscall"
	"time"

	"github.com/hashicorp/mdns"
	"github.com/manifoldco/promptui"
)

func RunSender() {
	done := make(chan struct{})

	// Make a channel for results and start listening
	entriesCh := make(chan *mdns.ServiceEntry, 4)
	go func() {
		timeout := time.After(time.Second * 3)
		entries := make([]*mdns.ServiceEntry, 0)
		hostnames := make([]string, 0)
	loop:
		for {
			select {
			case entry := <-entriesCh:
				entries = append(entries, entry)
				hostnames = append(hostnames, entry.Name)
			case <-timeout:
				break loop
			}
		}

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
		if err = SendTo(fmt.Sprintf("%s:%d", entry.Addr.String(), entry.Port)); err != nil {
			fmt.Println("Failed in sending the data to the receiver")
		}

		done <- struct{}{}
	}()

	// Start the lookup
	mdns.Lookup(ServiceName, entriesCh)
	<-done
}

func SendTo(address string) error {
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
