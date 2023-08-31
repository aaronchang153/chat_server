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

func (c *ChatServer) wsEndpoint(ctx *gin.Context) {
	websocketConn, err := c.upgrader.Upgrade(
		ctx.Writer,
		ctx.Request,
		nil,
	)
	if err != nil {
		c.logger.Println(err)
		return
	}

	HandleClient(websocketConn)
}
