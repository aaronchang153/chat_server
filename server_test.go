package main

import (
	"context"
	"log"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"nhooyr.io/websocket"
)

func TestServer(t *testing.T) {
	//t.Parallel()
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)

	s := httptest.NewServer(ChatServer{
		logger: logger,
	})
	defer s.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	conn, _, err := websocket.Dial(ctx, s.URL, &websocket.DialOptions{})
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close(websocket.StatusInternalError, "Failed to connect to chat server")

	for i := 0; i < 10; i++ {
		data := []byte("Hello world!")
		conn.Write(ctx, OPC_TEXT, data)

		_, recv, _ := conn.Read(ctx)
		log.Println(string(recv))
	}

	conn.Close(websocket.StatusNormalClosure, "")
}
