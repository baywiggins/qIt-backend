package server

import (
	"net/http"

	"github.com/baywiggins/qIt-backend/internal/config"
)

func StartServer() {
	// Add our API routes to the http handler
	HandleRoutes()
	// Set our http server to listen and serve
	http.ListenAndServe(config.API_URL, nil)
}