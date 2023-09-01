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

	connections := make([]*websocket.Conn, 10)
	for i := 0; i < 10; i++ {
		conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		if err != nil {
			log.Fatal("dial:", err)
		}
		defer CloseConn(conn)

		connections[i] = conn
	}

	for _, conn := range connections {
		conn.WriteMessage(websocket.TextMessage, []byte("Hello World!"))

		_, buffer, err := conn.ReadMessage()
		if err != nil {
			log.Fatal("read:", err)
		}
		//logger.Println(string(buffer))
		if string(buffer) != "Hello World!" {
			logger.Fatal("Mismatch on echo operation")
		}
	}
}

func CloseConn(c *websocket.Conn) {
	err := c.WriteMessage(websocket.CloseMessage, nil)
	if err != nil {
		logger.Println("Warning: Error while sending close message:", err)
	}
	c.Close()
}
