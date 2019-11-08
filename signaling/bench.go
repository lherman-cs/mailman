package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

func main() {
	n := 8000
	dial := func(name string) {
		u := url.URL{Scheme: "ws", Host: "localhost:4001", Path: "/ws"}

		c, _, err := websocket.DefaultDialer.Dial(u.String(), http.Header(map[string][]string{
			"name": []string{name},
		}))
		if err != nil {
			log.Fatal("dial:", err)
		}
		defer c.Close()

		select {}
	}

	for i := 0; i < n; i++ {
		go func(i int) {
			dial(fmt.Sprint(i))
		}(i)
		time.Sleep(time.Millisecond)
	}

	select {}
}
