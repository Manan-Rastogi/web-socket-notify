// This is the ENTRYPOINT of our Project

package main

import (
	"fmt"
	"log"
	"net/http"

	"websocket-notify/internal/websocket"
)

func main() {
	http.HandleFunc("/ws", websocket.HandleWebSocket)

	fmt.Println("🔌 WebSocket server running on http://localhost:8080/ws")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Server failed:", err)
	}
}
