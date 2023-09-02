package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const PORT_NUM string = ":8080"

type ChatServer struct {
	logger        *log.Logger
	upgrader      websocket.Upgrader
	clientMap     map[string]*ChatClientHandler
	notifRecvChan chan string
}

type ChatClientHandler struct {
	id            string
	logger        *log.Logger
	conn          *websocket.Conn
	notifSendChan chan string
	msgQ          MsgQueue
}

func RunServer(logger *log.Logger) error {
	ginEngine := gin.Default()

	upgrader := websocket.Upgrader{}
	server := ChatServer{
		logger:        logger,
		upgrader:      upgrader,
		clientMap:     make(map[string]*ChatClientHandler),
		notifRecvChan: make(chan string),
	}

	go server.broadcastHandler()

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

func (s *ChatServer) broadcastHandler() {
	for id := range s.notifRecvChan {
		message := s.clientMap[id].msgQ.Pop() //TODO: need to lock this queue (and maybe the clientMap too)

		for _, client := range s.clientMap {
			client.FowardMessageToClient(message)
		}
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
		id:            uuid.New().String(),
		logger:        s.logger,
		conn:          websocketConn,
		notifSendChan: s.notifRecvChan,
		msgQ:          NewMsgQueue(),
	}
	s.clientMap[handler.id] = &handler

	go handler.HandleClient()
}

func (c *ChatClientHandler) HandleClient() {
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			c.conn.Close()
			break
		}

		c.msgQ.Push(string(message))
		c.notifSendChan <- c.id
	}
}

func (c *ChatClientHandler) FowardMessageToClient(m string) {
	c.conn.WriteMessage(websocket.TextMessage, []byte(m))
}
