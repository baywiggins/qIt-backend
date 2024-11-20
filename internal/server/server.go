package server

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/baywiggins/qIt-backend/internal/config"
)

func StartServer(db *sql.DB) error {
	var err error;
	// Add our API routes to the http handler
	HandleRoutes(db)
	// Set our http server to listen and serve
	if err := http.ListenAndServe(config.API_URL, nil); err != nil {
		log.Fatalf("server failed: %s", err)
	}
	return err
}