package handlers

import (
    "encoding/json"
    "log"
    "net/http"
    "sync"

    "github.com/gorilla/websocket"
)

var (
    upgrader = websocket.Upgrader{
        ReadBufferSize:  1024,
        WriteBufferSize: 1024,
        CheckOrigin: func(r *http.Request) bool {
            return true // Allow all origins in development
        },
    }
    clients = make(map[*websocket.Conn]bool)
    clientsMutex sync.RWMutex
)

// Message represents a WebSocket message
type Message struct {
    Type    string      `json:"type"`
    Payload interface{} `json:"payload"`
}

// HandleWebSocket manages WebSocket connections and handles real-time communication.
// It upgrades HTTP connections to WebSocket protocol, manages client connections,
// and processes incoming messages.
func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
    // Validate the WebSocket upgrade request
    if !websocket.IsWebSocketUpgrade(r) {
        log.Printf("Not a WebSocket handshake")
        http.Error(w, "Expected WebSocket Upgrade", http.StatusBadRequest)
        return
    }

    // Upgrade the HTTP connection to WebSocket protocol
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Printf("Error upgrading connection: %v", err)
        return
    }
    defer conn.Close()

    // Set read/write deadlines and buffer sizes
    // 512KB max message size to prevent memory exhaustion
    conn.SetReadLimit(512 * 1024)
    // No timeout for long-lived connections
    conn.SetReadDeadline(time.Time{})

    // Register new client in the global clients map
    clientsMutex.Lock()
    clients[conn] = true
    clientsMutex.Unlock()

    // Ensure client is removed from the map when connection closes
    defer func() {
        clientsMutex.Lock()
        delete(clients, conn)
        clientsMutex.Unlock()
    }()

    // Main message processing loop
    // Continuously reads messages from the client until connection closes
    for {
        _, message, err := conn.ReadMessage()
        if err != nil {
            log.Printf("Error reading message: %v", err)
            break
        }

        // Unmarshal JSON message into Message struct
        var msg Message
        if err := json.Unmarshal(message, &msg); err != nil {
            log.Printf("Error unmarshaling message: %v", err)
            continue
        }

        // Process different message types
        switch msg.Type {
        case "ping":
            // Respond to ping messages with pong to maintain connection
            if err := conn.WriteJSON(Message{Type: "pong"}); err != nil {
                log.Printf("Error sending pong: %v", err)
            }
        }
    }
}

// BroadcastTodoUpdate sends a todo update to all connected clients
func BroadcastTodoUpdate(todo interface{}, messageType string) {
    message := Message{
        Type:    messageType,
        Payload: todo,
    }

    clientsMutex.RLock()
    defer clientsMutex.RUnlock()

    for client := range clients {
        if err := client.WriteJSON(message); err != nil {
            log.Printf("Error broadcasting message: %v", err)
            client.Close()
            delete(clients, client)
        }
    }
}