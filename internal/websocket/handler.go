package websocket

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // allow all origins for dev
	},
}

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	deviceID := r.URL.Query().Get("deviceID")
	if deviceID == "" {
		http.Error(w, "deviceID is required", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade failed:", err)
		return
	}
	defer conn.Close()

	log.Printf("✅ Connected: %s\n", deviceID)

	for {
		mt, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("❌ Disconnected: %s (%v)\n", deviceID, err)
			break
		}
		log.Printf("📩 Message from %s: %s", deviceID, msg)

		// Echo back the message
		err = conn.WriteMessage(mt, msg)
		if err != nil {
			log.Println("Write error:", err)
			break
		}
	}
}
