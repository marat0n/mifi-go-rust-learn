package main

import "fmt"

func main() {
	// var myVar int = 0

	fmt.Println(multipleStrings("Hello ", 11))
}

func multipleStrings(str string, num uint) string {
	if num == 0 {
		return ""
	}

	resStr := ""
	for range num {
		resStr = resStr + str
	}
	return resStr
}


// Scrapper
// 1. HTTP request
// 2. JS execution
// 3. Virtual Browser Core | Conteinerized Chromium
// 4. JS -> NodeJS -> Scrapping Framework


// Back-end (REST API, gRPC, GraphQL API)


// Chat | Global chat
// WebSocket API
// Real-time communication
// SQLite database (messages + users)

// UNIX Philosophy
// Vi (visual edit) -> Vim (vi modern) -> Nvim (neo vi modern)
