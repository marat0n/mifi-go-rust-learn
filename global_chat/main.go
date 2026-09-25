package main

import (
	"fmt"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Message []byte

// State structure for every websocket connection.
// It holds all clients' connections and other usefull fields for dynamic messages broadcasting.
type State struct {
	mutex    sync.Mutex
	clients  map[*websocket.Conn]bool
	messages []Message
}

func main() {
	// Creating a default HTTP router (using the Gin library)
	router := gin.Default()

	// Creating an instance of our websocket state
	state := State{
		// Allocating a memory for websocket connections
		// (a dictionary assosiation of connection -> online status)
		clients: make(map[*websocket.Conn]bool),
	}

	// A "Ping-Pong" method for testing our back-end
	router.GET("/ping", func(ctx *gin.Context) {
		ctx.Writer.WriteString("Pong")
	})

	// Method for websocket connection
	router.GET("/chat", state.chat)

	// Running all router's methods
	router.Run()
}

// A default websocket upgrader.
// WebSocket is a protocol based on the HTTP protocol,
// 	so you need to upgrade the connection once it executed the default "GET" method.
var upgrader = websocket.Upgrader{}

func (state *State) chat(ctx *gin.Context) {
	// Getting a "Writer" and "Request" models from the http-context
	writer, request := ctx.Writer, ctx.Request

	// Trying to upgrade the connection (from default GET HTTP to the WebSocket)
	conn, err := upgrader.Upgrade(writer, request, nil)

	// Checking the error after connection upgrading
	if err != nil {
		_, err = writer.WriteString("Couldn't upgrade the connection... :(")
		if err != nil {
			fmt.Printf("Error occurred when tried to upgrade the connection: %s\n", err)
		}
		return
	}

	// The `defer` keyword is used to run a command as the last thing in this function
	// 	so we are closing the connection after everything else is done.
	defer conn.Close()

	// Sending all saved messages from the `state` to a new connected client
	for _, msg := range state.messages {
		conn.WriteMessage(websocket.TextMessage, msg)
	}

	// We are trying to lock the mutex.
	// If this mutex was locked by another connection, we must to wait untill it unlocked.
	// When it'll be unlocked we will lock ourselves, so no other connections could interrupt us here.
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

		state.mutex.Lock()
		state.messages = append(state.messages, message)
		state.mutex.Unlock()

		fmt.Printf("Recieved msg: %s\n", message)
		go state.broadcast(message)
	}
}

func (state *State) broadcast(message []byte) {
	state.mutex.Lock()
	defer state.mutex.Unlock()

	for client := range state.clients {
		if err := client.WriteMessage(websocket.TextMessage, message); err != nil {
			fmt.Printf("write error: %s\n", err)
			client.Close()
			delete(state.clients, client)
		}
	}
}
