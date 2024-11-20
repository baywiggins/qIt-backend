package server

import (
	"database/sql"

	"github.com/baywiggins/qIt-backend/internal/api/handlers"
)

func HandleRoutes(db *sql.DB) {
	handlers.HandleSpotifyAuthRoutes(db)
	handlers.HandleSpotifyControllerRoutes(db)
	handlers.HandleAuthRoutes(db)
}