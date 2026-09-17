package main

import (
	"fmt"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type State struct {
	mutex    sync.Mutex
	clients  map[*websocket.Conn]bool
}

func main() {
	router := gin.Default()

	state := State{
		clients: make(map[*websocket.Conn]bool),
	}

	router.GET("/ping", func(ctx *gin.Context) {
		ctx.Writer.WriteString("Pong")
	})
	router.GET("/chat", state.chat)

	router.Run()
}

var upgrader = websocket.Upgrader{}

func (state *State) chat(ctx *gin.Context) {
	writer, request := ctx.Writer, ctx.Request

	conn, err := upgrader.Upgrade(writer, request, nil)
	if err != nil {
		_, err = writer.WriteString("Couldn't upgrade the connection... :(")
		if err != nil {
			fmt.Printf("Error occurred when tried to upgrade the connection: %s\n", err)
		}
		return
	}
	defer conn.Close()

	state.mutex.Lock()
	state.clients[conn] = true
	state.mutex.Unlock()

	defer func() {
		state.mutex.Lock()
		delete(state.clients, conn)
		state.mutex.Unlock()
	}()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Printf("Error on reading a message: %s\n", err)
			break
		}

		fmt.Printf("Recieved msg: %s\n", message)
		go state.broadcast(message)
	}
}

func (state *State) broadcast(message []byte) {
	// TODO: Заполнить эту функцию, чтобы отправить всем пользователям (state.clients) сообщение msg
}
