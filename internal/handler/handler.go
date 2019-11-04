package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
	msg "github.com/lherman-cs/mailman/internal/message"
	"github.com/lherman-cs/mailman/internal/peer"
	"go.uber.org/zap"
)

var upgrader = websocket.Upgrader{} // use default options

type Handler struct {
	http.Handler
	Params
	rl atomic.Value
	wl atomic.Value
}

type Params struct {
	Sender     bool
	Port       int
	Log        *zap.SugaredLogger
	DataInput  io.Reader
	DataOutput io.Writer
}

func New(params Params) *Handler {
	if params.DataInput == nil {
		params.DataInput = os.Stdin
	}

	if params.DataOutput == nil {
		params.DataOutput = os.Stdout
	}

	a := Handler{Params: params}

	r := chi.NewRouter()
	// setup middlewares
	// debugger
	r.Use(func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			a.Log.Debugw("started",
				"endpoint", r.URL.Path,
			)

			start := time.Now()
			defer a.Log.Debugw("finished",
				"endpoint", r.URL.Path,
				"duration", time.Now().Sub(start),
			)
			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	})

	// recoverer
	r.Use(func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				rvr := recover()
				if rvr == nil {
					return
				}

				status := http.StatusInternalServerError
				err, ok := rvr.(*HTTPError)
				if ok {
					status = err.status
				}

				reason := fmt.Sprint(err)

				http.Error(w, reason, status)
				a.Log.Errorw(reason,
					"status", status,
					"path", r.URL.Path,
				)
			}()

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	})

	r.Post("/log", a.log)
	r.Get("/getHost", a.getHost)
	r.Get("/getPeerType", a.getPeerType)
	r.HandleFunc("/connect", a.connect)
	a.Handler = r

	return &a
}

func (a *Handler) log(w http.ResponseWriter, r *http.Request) {
	type LogFormat struct {
		Level   string                 `json:"level"`
		Message string                 `json:"message"`
		Fields  map[string]interface{} `json:"fields"`
	}

	var l LogFormat
	err := json.NewDecoder(r.Body).Decode(&l)
	if err != nil {
		panic(&HTTPError{
			status: http.StatusBadRequest,
			reason: "log's body has to be in valid json format",
		})
	}

	fields := make([]interface{}, len(l.Fields)*2)
	for key, val := range l.Fields {
		fields = append(fields, key, val)
	}

	switch l.Level {
	case "debug":
		a.Log.Debugw(l.Message, fields...)
	case "warn":
		a.Log.Warnw(l.Message, fields...)
	case "error":
		a.Log.Errorw(l.Message, fields...)
	default:
		a.Log.Infow(l.Message, fields...)
	}
}

func (a *Handler) getHost(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "ws://localhost:%d", a.Port)
	if err != nil {
		panic(err)
	}
}

func (a *Handler) getPeerType(w http.ResponseWriter, r *http.Request) {
	var err error
	if a.Sender {
		_, err = fmt.Fprint(w, "sender")
	} else {
		_, err = fmt.Fprint(w, "receiver")
	}

	if err != nil {
		panic(err)
	}
}

func (a *Handler) connect(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		panic(err)
	}

	cSync := newConn(c)
	p, err := peer.NewPeer(context.Background(), cSync, a.Log)
	if err != nil {
		panic(err)
	}
	defer p.Close()

	if a.Sender {
		err = a.senderHandler(cSync, p)
	} else {
		err = a.receiverHandler(cSync, p)
	}

	if err != nil {
		panic(err)
	}
}

func (a *Handler) senderHandler(conn *connSync, p *peer.Peer) error {
	a.Log.Info("handling sender")

	conn.WriteJSON(msg.StateMessage("Connecting through peer-to-peer"))
	writer, err := p.Connect()
	if err != nil {
		return err
	}
	defer writer.Close()

	conn.WriteJSON(msg.StateMessage("Sending data"))
	_, err = io.Copy(writer, a.DataInput)
	conn.WriteJSON(msg.StateMessage("Data has been sent"))
	return err
}

func (a *Handler) receiverHandler(conn *connSync, p *peer.Peer) error {
	a.Log.Info("handling receiver")

	conn.WriteJSON(msg.StateMessage("Connecting through peer-to-peer"))
	reader, err := p.Wait()
	if err != nil {
		return err
	}
	defer reader.Close()

	conn.WriteJSON(msg.StateMessage("Receiving data"))
	_, err = io.Copy(a.DataOutput, reader)
	conn.WriteJSON(msg.StateMessage("Data has been received"))
	return err
}
