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
		messageType, message, err := c.conn.ReadMessage()
		if err != nil {
			c.logger.Println(err)
			break
		}

		if messageType == websocket.CloseMessage {
			break
		}
		c.conn.WriteMessage(messageType, message)
	}
}
