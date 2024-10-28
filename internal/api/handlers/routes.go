package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

func HandleRoutes() {
	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		log.Println("'/status' was called")
		w.Header().Set("Content-Type", "application/json")

		status := map[string]string{"status": "OK"}
		json.NewEncoder(w).Encode(status)
	})
}