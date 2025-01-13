package pubsub

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/baywiggins/qIt-backend/internal/config"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

// Set up redis client
var redisClient = redis.NewClient(&redis.Options{
	Addr: config.RedisAddress,
})

// Websocket upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {return true},
}

func handleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
		return
	}
	defer conn.Close()

	roomID := r.URL.Query().Get("room")
	if roomID == "" {
		conn.WriteMessage(websocket.TextMessage, []byte("Room ID is required"))
		return
	}

	ctx := context.Background()
	pubsub := redisClient.Subscribe(ctx, "room-"+roomID)
	defer pubsub.Close()

	log.Println("User connected to room:", roomID)

	// Periodic ping to keep the connection alive
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	go func() {
		for range ticker.C {
			if err := conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				log.Println("Ping error:", err)
				return
			}
		}
	}()

	// Listen to messages from Redis pub/sub
	for msg := range pubsub.Channel() {
		if err := conn.WriteMessage(websocket.TextMessage, []byte(msg.Payload)); err != nil {
			log.Println("Write error:", err)
			return
		}
	}
}


func RunWebSocket() {
	http.HandleFunc("/votes", handleConnection)
	log.Println("WebSocket server started on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}