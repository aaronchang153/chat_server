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

var addr = flag.String("addr", "localhost:8080", "http service address")

func TestServer(t *testing.T) {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)

	go RunServer(logger)

	time.Sleep(3 * time.Second)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "websocket"}

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer conn.Close()

	conn.WriteMessage(websocket.TextMessage, []byte("Hello World!"))

	_, buffer, err := conn.ReadMessage()
	if err != nil {
		log.Fatal("read:", err)
	}
	logger.Println(string(buffer))
}
