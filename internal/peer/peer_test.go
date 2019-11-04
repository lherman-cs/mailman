package peer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"sync"
	"testing"

	"go.uber.org/zap"
)

type MockConn struct {
	in  <-chan *bytes.Buffer
	out chan<- *bytes.Buffer
}

var _ Conn = &MockConn{}

func NewMockConns() (conn1, conn2 *MockConn) {
	chan1 := make(chan *bytes.Buffer, 128)
	chan2 := make(chan *bytes.Buffer, 128)

	conn1 = &MockConn{in: chan1, out: chan2}
	conn2 = &MockConn{in: chan2, out: chan1}

	return
}

func (c *MockConn) ReadJSON(v interface{}) error {
	msg, ok := <-c.in
	if !ok {
		return nil
	}
	return json.NewDecoder(msg).Decode(v)
}

func (c *MockConn) WriteJSON(v interface{}) error {
	var msg bytes.Buffer
	err := json.NewEncoder(&msg).Encode(v)
	if err != nil {
		return err
	}

	c.out <- &msg
	return nil
}

func (c *MockConn) Close() error {
	close(c.out)
	return nil
}

func TestMockConn(t *testing.T) {
	type Data struct {
		Message string
	}
	conn1, conn2 := NewMockConns()

	err := conn1.WriteJSON(&Data{Message: "echo1"})
	if err != nil {
		t.Fatal("failed in writing json")
	}

	err = conn1.WriteJSON(&Data{Message: "echo2"})
	if err != nil {
		t.Fatal("failed in writing json")
	}

	var d Data
	err = conn2.ReadJSON(&d)
	if err != nil {
		t.Fatal("failed in reading json")
	}

	if d.Message != "echo1" {
		t.Fatal("message is invalid")
	}

	err = conn2.ReadJSON(&d)
	if err != nil {
		t.Fatal("failed in reading json")
	}

	if d.Message != "echo2" {
		t.Fatal("message is invalid")
	}
}

func TestEcho(t *testing.T) {
	data := "echo"
	var wg sync.WaitGroup
	wg.Add(2)
	conn1, conn2 := NewMockConns()

	// receiver
	go func() {
		defer wg.Done()
		l, err := zap.NewDevelopment()
		if err != nil {
			t.Fatal(err)
		}
		peer, err := NewPeer(context.Background(), conn1, l.Sugar())
		if err != nil {
			t.Fatalf("failed in initializing receiver: %+v", err)
		}
		defer peer.Close()
		log.Println("receiver has been created")

		reader, err := peer.Wait()
		if err != nil {
			t.Fatalf("failed in waiting: %+v", err)
		}
		defer reader.Close()
		log.Println("receiver got a sender")

		var out bytes.Buffer
		_, err = io.Copy(&out, reader)
		if err != nil {
			t.Fatalf("failed in getting data from sender: %+v", err)
		}
		log.Println("receiver got data")

		if out.String() != data {
			t.Fatal("data received is invalid")
		}

		log.Println("got", data)
	}()

	// sender
	go func() {
		// time.Sleep(time.Second)
		defer wg.Done()
		l, err := zap.NewDevelopment()
		if err != nil {
			t.Fatal(err)
		}
		peer, err := NewPeer(context.Background(), conn2, l.Sugar())
		if err != nil {
			t.Fatalf("failed in initializing sender: %+v", err)
		}
		log.Println("sender has been created")
		defer peer.Close()

		writer, err := peer.Connect()
		if err != nil {
			t.Fatalf("failed in connecting: %+v", err)
		}
		defer writer.Close()
		log.Println("sender is connected")

		_, err = fmt.Fprint(writer, data)
		if err != nil {
			t.Fatalf("failed in sending data to the receiver: %+v", err)
		}
		log.Println("sender sent data")

	}()

	wg.Wait()
}
