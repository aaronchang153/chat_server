package main

import (
	"context"
	"log"
	"net/http"

	"nhooyr.io/websocket"
)

const (
	OPC_CONT websocket.MessageType = 0x0
	OPC_TEXT websocket.MessageType = 0x1
	OPC_BIN  websocket.MessageType = 0x2
)

type ChatServer struct {
	logger *log.Logger
}

func (s ChatServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, nil) //Can specify and check supported protocols if needed
	if err != nil {
		s.logger.Printf("%v\n", err)
		return
	}
	defer conn.Close(websocket.StatusInternalError, "Internal error")

	context := r.Context()

	var msgType websocket.MessageType
	var data []byte
	for {
		msgType, data, err = conn.Read(context)
		if err != nil {
			s.logger.Printf("Error encountered while reading from WebSocket: %v\n", err)
			break
		}

		switch msgType {
		case OPC_TEXT:
			err = s.HandleText(conn, &data, context)
		default:
			s.logger.Printf("Ignoring WebSocket message with opcode %#x\n", msgType)
		}

		if err != nil {
			s.logger.Printf("Error encountered while handling last message: %v\n", err)
		}
	}
}

func (s ChatServer) HandleText(conn *websocket.Conn, data *[]byte, ctx context.Context) error {
	return conn.Write(ctx, OPC_TEXT, *data)
}
