package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/baywiggins/qIt-backend/internal/api/middlewares"
	"github.com/baywiggins/qIt-backend/internal/models"
	"github.com/baywiggins/qIt-backend/pkg/utils"

	"github.com/google/uuid"
)

// Obviously change this to more secure storage in the future
var users = map[string]string{}

type Handler struct {
	DB *sql.DB
}

type UserData struct {
	Username string `json:"username"`
	Password string `json:"password"`
	UserState string `json:"state"`
}

func HandleAuthRoutes(db *sql.DB) {
	handler := &Handler{DB: db}
	// Handle our GET endpoint to create an account
	http.Handle("POST /account/create", middlewares.LoggingMiddleware(http.HandlerFunc(handler.handleCreateAccount)))
	// Handle our GET endpoint to login
	http.Handle("POST /account/login", middlewares.LoggingMiddleware(http.HandlerFunc(handler.handleLogin)))
}

func (h *Handler) handleCreateAccount(w http.ResponseWriter, r *http.Request) {
	var err error;

	var userData UserData;

	err = json.NewDecoder(r.Body).Decode(&userData)
	if err != nil {
		log.Printf("ERROR: %s \n", err)
		utils.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error in handleCreateAccount: '%s' \n", err.Error()))
		return
	}

	userID := uuid.NewString()
	hashedPassword, err := utils.HashPassword(userData.Password)
	if err != nil {
		log.Printf("ERROR: %s \n", err)
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error in handleCreateAccount: '%s' \n", err.Error()))
		return
	}
	// Add user to database
	newUser := models.User {
		ID: userID, 
		Username: userData.Username, 
		Password: hashedPassword, 
		UserState: userData.UserState,
	}

	if err := models.InsertUser(h.DB, newUser); err != nil {
		log.Printf("ERROR: %s \n", err)
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error in handleCreateAccount: '%s' \n", err.Error()))
		return
	}

	token, err := utils.GenerateJWTToken(userID)
	if err != nil {
		log.Printf("ERROR: %s \n", err)
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error in handleCreateAccount: '%s' \n", err.Error()))
		return
	}
	w.Write([]byte(token))
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	var err error;

	var loginData UserData;

	err = json.NewDecoder(r.Body).Decode(&loginData)
	if err != nil {
		log.Printf("ERROR: %s \n", err)
		utils.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error in handleLogin: '%s' \n", err.Error()))
		return
	}

	// This definitely needs an update too 💀
	// Also password checking xd
	for userID, username := range users {
		if username == loginData.Username {
			token, err := utils.GenerateJWTToken(userID)
			if err != nil {
				log.Printf("ERROR: %s \n", err)
				utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error in handleLogin: '%s' \n", err.Error()))
				return
			}
			w.Write([]byte(token))
			return
		}
	}
	log.Printf("ERROR: %s \n", "invalid credentials")
	utils.RespondWithError(w, http.StatusUnauthorized, fmt.Sprintf("Error in handleLogin: '%s' \n", "invalid credentials"))
}