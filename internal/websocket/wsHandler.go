package websocket

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	clients = make(map[string]*websocket.Conn) // deviceId -> websocket.Conn
)

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// 1. Extract deviceId from URL path like
	deviceID := r.URL.Query().Get("deviceID")
	if deviceID == "" {
		http.Error(w, "deviceID is required", http.StatusBadRequest)
		return
	}

	// 2. Upgrade connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade failed:", err)
		return
	}
	defer conn.Close()

	// 3. Register client
	clients[deviceID] = conn
	log.Printf("✅ : %s Registered.\n", deviceID)

	// 4. Simulate a server push to just this device after 5 seconds
	go func(id string) {
		time.Sleep(5 * time.Second)
		msg := fmt.Sprintf("Hello %s! This is a server-triggered message.", id)

		if clientConn, ok := clients[id]; ok {
			err := clientConn.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				log.Printf("[ERROR] Writing to %s failed: %v", id, err)
			} else {
				log.Printf("[PUSHED] Sent to %s", id)
			}
		}
	}(deviceID)

	// 5. Read messages (to keep connection open)
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("❌ Disconnected: %s (%v)\n", deviceID, err)
			break
		}
		log.Printf("[RECV] %s: %s", deviceID, string(msg))

	}
}
