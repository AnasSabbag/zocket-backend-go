package handlers

import (
	"net/http"
	"sync"
	"fmt"
	"github.com/AnasSabbag/task-manager/models"
	"github.com/gorilla/websocket"
)

// WebSocket Upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan models.Message)
var mutex = &sync.Mutex{}

// WebSocketHandler handles WebSocket connections

func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("WebSocket Upgrade Error:", err)
		http.Error(w, "Could not open WebSocket connection", http.StatusBadRequest)
		return
	}
	defer conn.Close()
	fmt.Println("WebSocket connection established")

	for {
		// Just keep reading or writing to keep it open
		_, _, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("WebSocket closed:", err)
			break
		}
	}
}



// StartWebSocketListener continuously sends broadcast messages to all clients
func StartWebSocketListener() {
	for {
		msg := <-broadcast
		mutex.Lock()
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				client.Close()
				delete(clients, client)
			}
		}
		mutex.Unlock()
	}
}

// SendMessageToClients broadcasts a message to all clients
func SendMessageToClients(msg models.Message) {
	broadcast <- msg
}
