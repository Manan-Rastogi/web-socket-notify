package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// Payload structure
type PushPayload struct {
	DeviceID string `json:"deviceId"`
	Message  string `json:"message"`
}

var mu sync.RWMutex // for thread-safe access

// HTTP handler to push messages
func PushHandler(w http.ResponseWriter, r *http.Request) {
	var payload PushPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	mu.RLock()
	conn, ok := clients[payload.DeviceID]
	mu.RUnlock()

	if !ok {
		http.Error(w, "Device not connected", http.StatusNotFound)
		return
	}

	// Send message via WebSocket
	err := conn.WriteMessage(websocket.TextMessage, []byte(payload.Message))
	if err != nil {
		log.Printf("❌ Error sending to %s: %v\n", payload.DeviceID, err)
		http.Error(w, "Failed to send", http.StatusInternalServerError)
		return
	}

	log.Printf("✅ Sent to %s: %s\n", payload.DeviceID, payload.Message)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Message sent"))
}
