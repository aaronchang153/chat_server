package main

import (
	"flag"
	"log"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

var logger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)

var addr = flag.String("addr", "localhost:8080", "http service address")

func TestServer(t *testing.T) {
	t.Parallel()
	go RunServer(logger)
	time.Sleep(1 * time.Second)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "websocket"}

	const NUM_CONNECTIONS int = 10
	connections := make([]*websocket.Conn, NUM_CONNECTIONS)
	for i := 0; i < NUM_CONNECTIONS; i++ {
		conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		if err != nil {
			log.Fatal("dial:", err)
		}
		defer CloseConn(conn)

		connections[i] = conn
	}

	for _, conn := range connections {
		msg := "Hello World!"
		conn.WriteMessage(websocket.TextMessage, []byte(msg))

		conn.SetReadDeadline(time.Now().Add(3 * time.Second))

		var resp []byte
		_, resp, err := conn.ReadMessage()
		if err != nil {
			log.Fatal("read:", err)
		}
		if string(resp) != msg {
			logger.Fatal("Mismatch on echo operation")
		}
	}
}

func CloseConn(c *websocket.Conn) {
	err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		logger.Println("Warning: Error while sending close message:", err)
	}
}
