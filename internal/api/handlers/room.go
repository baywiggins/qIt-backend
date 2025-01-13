package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/baywiggins/qIt-backend/internal/api/middlewares"
	"github.com/baywiggins/qIt-backend/internal/config"
	"github.com/baywiggins/qIt-backend/pkg/utils"
	"github.com/redis/go-redis/v9"
)

var redisClient = redis.NewClient(&redis.Options{
	Addr: config.RedisAddress,
})

func publishToRoom(ctx context.Context, roomID string, message []byte) (error) {
	var err error;
	err = redisClient.Publish(ctx, "room-"+roomID, message).Err()
	
	return err
}

type VoteRequest struct {
	RoomID string `json:"room_id"`
	UserID string `json:"user_id"`
	Vote   bool   `json:"vote"`
}

type VoteUpdate struct {
    RoomID    string            `json:"roomID"`
    VoteCount map[string]int    `json:"voteCount"`
}

func HandleRoomRoutes(db *sql.DB) {
	handler := &Handler{DB: db}
	// Create room
	http.Handle("/room/create", middlewares.LoggingMiddleware(middlewares.AuthMiddleware(http.HandlerFunc(handler.handleCreateRoom))))

	// Join? room

	// Handle voting
	http.Handle("/room/vote", middlewares.LoggingMiddleware(http.HandlerFunc(handler.handleVote)))
}

func (h *Handler) handleCreateRoom(w http.ResponseWriter, r *http.Request){}

func (h *Handler) handleVote(w http.ResponseWriter, r *http.Request) {
	var voteReq VoteRequest;
	if err := json.NewDecoder(r.Body).Decode(&voteReq); err != nil {
		log.Printf("Error in handleVote: '%s'", err.Error())
		utils.RespondWithError(w, http.StatusBadRequest, "Missing required info")
		return
	}

	// Simulate vote processing (you'd update a database in real logic)
    updatedVotes := map[string]int{
        "OptionA": 5,
        "OptionB": 3,
    }

	// Publish to room channel
	ctx := context.Background()
	update := VoteUpdate{
		RoomID: voteReq.RoomID,
		VoteCount: updatedVotes,
	}
	
	message, _ := json.Marshal(update)

	if err := publishToRoom(ctx, voteReq.RoomID, message); err != nil {
		log.Printf("Error in handleVote: '%s'", err.Error())
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to broadcast new votes")
		return
	}
}