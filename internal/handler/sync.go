package handler

import (
	"sync"

	"github.com/gorilla/websocket"
)

type connSync struct {
	*websocket.Conn
	l sync.Mutex
}

func newConn(conn *websocket.Conn) *connSync {
	return &connSync{Conn: conn}
}

func (c *connSync) WriteJSON(v interface{}) error {
	c.l.Lock()
	defer c.l.Unlock()

	return c.Conn.WriteJSON(v)
}
