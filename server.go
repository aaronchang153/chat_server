package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const PORT_NUM string = ":8080"

type ChatServer struct {
	logger   *log.Logger
	upgrader websocket.Upgrader
}

type ChatClientHandler struct {
	logger *log.Logger
	conn   *websocket.Conn
}

type ClientMessage struct {
	MessageType int
	Message     string
}

const (
	MSG_ERROR int = -1
	MSG_ECHO  int = 0
)

func RunServer(logger *log.Logger) error {
	ginEngine := gin.Default()

	upgrader := websocket.Upgrader{}
	server := ChatServer{
		logger:   logger,
		upgrader: upgrader,
	}

	ginEngine.GET("/websocket", server.wsEndpoint)
	return ginEngine.Run(PORT_NUM)
}

func HandleClient(conn *websocket.Conn) {
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			break
		}

		conn.WriteMessage(messageType, message)
	}
}

func (s *ChatServer) wsEndpoint(ctx *gin.Context) {
	websocketConn, err := s.upgrader.Upgrade(
		ctx.Writer,
		ctx.Request,
		nil,
	)
	if err != nil {
		s.logger.Println(err)
		return
	}

	handler := ChatClientHandler{
		logger: s.logger,
		conn:   websocketConn,
	}
	go handler.HandleClient()
}

func (c *ChatClientHandler) HandleClient() {
	for {
		//messageType, message, err := c.conn.ReadMessage()
		var message ClientMessage
		err := c.conn.ReadJSON(&message)
		if err != nil {
			c.conn.Close()
			break
		}

		c.ProcessMessage(message)
	}
}

func (c *ChatClientHandler) ProcessMessage(m ClientMessage) error {
	switch m.MessageType {
	case MSG_ECHO:
		c.conn.WriteJSON(m)
	default:
		response := ClientMessage{
			MessageType: MSG_ERROR,
			Message:     "Malformatted request",
		}
		c.conn.WriteJSON(response)
	}
	return nil
}
