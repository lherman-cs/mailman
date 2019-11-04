//go:generate go run -tags generate gen.go
package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"

	"github.com/go-chi/chi"
	"github.com/lherman-cs/lorca"
	"github.com/lherman-cs/mailman/internal/handler"
)

func main() {
	defer logger.Sync()

	var params handler.Params
	fi, err := os.Stdin.Stat()
	if err != nil {
		logger.Fatal(err)
	}

	if fi.Mode()&os.ModeNamedPipe != 0 {
		params.Sender = true
	}

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		logger.Fatal(err)
	}

	params.Port = listener.Addr().(*net.TCPAddr).Port
	params.Log = logger
	h := handler.New(params)

	fs := http.FileServer(assets)
	r := chi.NewRouter()
	r.Mount("/", fs)
	r.Mount("/api", h)

	url := fmt.Sprintf("http://localhost:%d", params.Port)
	logger.Infof("listening on: %s", url)
	if !production {
		err = http.Serve(listener, r)
		if err != nil {
			logger.Fatal(err)
		}
		return
	}

	go http.Serve(listener, r)
	args := []string{"--start-fullscreen", "--disable-dev-tools"}
	ui, err := lorca.New("", userDir, 480, 320, logger, args...)
	if err != nil {
		logger.Fatal(err)
	}
	defer ui.Close()
	err = ui.Load(url)
	// if can't load google chrome, fallback to use any existing browser instead
	if err != nil {
		err = open(url)
		if err != nil {
			logger.Fatal(err)
		}
	}

	// Wait until the interrupt signal arrives or browser window is closed
	sigc := make(chan os.Signal)
	signal.Notify(sigc, os.Interrupt)
	select {
	case <-sigc:
	case <-ui.Done():
	}

	logger.Info("exiting...")
}
