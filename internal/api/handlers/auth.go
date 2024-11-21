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
	// Handle test auth endpoint to ensure users got a valid token
	http.Handle("GET /account/test-auth", middlewares.LoggingMiddleware(middlewares.AuthMiddleware(http.HandlerFunc(handleTestAuth))))
}

func (h *Handler) handleCreateAccount(w http.ResponseWriter, r *http.Request) {
	var err error;

	var userData UserData;

	err = json.NewDecoder(r.Body).Decode(&userData)
	if err != nil {
		log.Printf("ERROR in handleCreateAccount: %s \n", err)
		utils.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error in handleCreateAccount: '%s' \n", err.Error()))
		return
	}

	userID := uuid.NewString()
	hashedPassword, err := utils.HashPassword(userData.Password)
	if err != nil {
		log.Printf("ERROR in handleCreateAccount: %s \n", err)
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
		log.Printf("ERROR in handleCreateAccount: %s \n", err)
		if err.Error() == "UNIQUE constraint failed: Users.username" || err.Error() == "UNIQUE constraint failed: Users.user_state" {
			utils.RespondWithError(w, http.StatusBadRequest, "user already exists")
		return
		}
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error: '%s' \n", err.Error()))
		return
	}

	token, err := utils.GenerateJWTToken(userID)
	if err != nil {
		log.Printf("ERROR in handleCreateAccount: %s \n", err)
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error in handleCreateAccount: '%s' \n", err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	jsonEncoder := json.NewEncoder(w)
	jsonEncoder.Encode(map[string]string{
		"token": token,
		"uuid": userID,
	})
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	var err error;

	var loginData UserData;

	fmt.Println("reached 1")

	err = json.NewDecoder(r.Body).Decode(&loginData)
	if err != nil {
		log.Printf("ERROR in handleLogin: %s \n", err)
		utils.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error in handleLogin: '%s' \n", err.Error()))
		return
	}

	fmt.Println(loginData)

	// Get user data from the db
	username := loginData.Username
	user, err := models.GetByUserName(h.DB, username)
	if err != nil {
		log.Printf("ERROR in handleLogin: %s \n", err)
		statusCode := http.StatusInternalServerError
		if err.Error() == fmt.Sprintf("'%s' does not exist", username) {
			statusCode = http.StatusNotFound
		}
		utils.RespondWithError(w, statusCode, fmt.Sprintf("Error in handleLogin: '%s' \n", err.Error()))
		return
	}
	fmt.Println("reached 1")
	// Compare given password to hashed password stored in db
	matches := utils.DoPasswordsMatch(user.Password, loginData.Password)
	
	if (matches) {
		// Passwords match, login user
		token, err := utils.GenerateJWTToken(user.ID)
		if err != nil {
			log.Printf("ERROR in handleLogin: %s \n", err)
			utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error in handleCreateAccount: '%s' \n", err.Error()))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		jsonEncoder := json.NewEncoder(w)
		jsonEncoder.Encode(map[string]string{
			"token": token,
			"uuid": user.ID,
		})
	} else {
		log.Printf("ERROR user '%s' attempted to login with invalid credentials", username)
		utils.RespondWithError(w, http.StatusUnauthorized, "invalid credentials")
	}
}

func handleTestAuth(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("authorized!"))
}