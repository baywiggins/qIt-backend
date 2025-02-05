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
	// Handle refreshing JWT token
	http.Handle("POST /account/refresh-token", middlewares.LoggingMiddleware(http.HandlerFunc(handler.handleRefreshToken)))
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

	if userData.Username == "" || userData.Password == "" {
		log.Printf("ERROR in handleCreateAccount: %s \n", err)
		utils.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error in handleCreateAccount: '%s' \n", "missing username or password"))
		return
	}

	if userData.UserState == "" {
		log.Printf("ERROR in handleCreateAccount: %s \n", err)
		utils.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error in handleCreateAccount: '%s' \n", "missing state"))
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
			utils.RespondWithError(w, http.StatusBadRequest, "user or state already exists")
		return
		}
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error: '%s' \n", err.Error()))
		return
	}

	token, exp, err := utils.GenerateJWTToken(userID)
	if err != nil {
		log.Printf("ERROR in handleCreateAccount: %s \n", err)
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error in handleCreateAccount: '%s' \n", err.Error()))
		return
	}
	refreshToken, err := utils.GenerateRefreshToken(userID)
	if err != nil {
		log.Printf("ERROR in handleCreateAccount: %s \n", err)
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error in handleCreateAccount: '%s' \n", err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	jsonEncoder := json.NewEncoder(w)
	jsonEncoder.Encode(map[string]string{
		"token": token,
		"token_exp": exp,
		"refresh_token": refreshToken,
		"uuid": userID,
	})
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	var err error;

	var loginData UserData;

	err = json.NewDecoder(r.Body).Decode(&loginData)
	if err != nil {
		log.Printf("ERROR in handleLogin: %s \n", err)
		utils.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error in handleLogin: '%s' \n", err.Error()))
		return
	}


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

	// Ensure the user finished creating their account
	if (user.FinishedCreating == 0) {
		// They never finished, delete account and return 404
		err = models.DeleteUserByID(h.DB, user.ID)
		if err != nil {
			log.Printf("ERROR in handleLogin: %s \n", err)
			utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error in handleLogin: '%s' \n", err.Error()))
			return
		}
		log.Printf("ERROR in handleLogin: '%s' not found \n", user.Username)
		utils.RespondWithError(w, http.StatusNotFound, "Error: 'user not found'")
		return
	}

	// Compare given password to hashed password stored in db
	matches := utils.DoPasswordsMatch(user.Password, loginData.Password)
	
	if (matches) {
		// Passwords match, login user
		token, exp, err := utils.GenerateJWTToken(user.ID)
		if err != nil {
			log.Printf("ERROR in handleCreateAccount: %s \n", err)
			utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error in handleCreateAccount: '%s' \n", err.Error()))
			return
		}
		refreshToken, err := utils.GenerateRefreshToken(user.ID)
		if err != nil {
			log.Printf("ERROR in handleCreateAccount: %s \n", err)
			utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error in handleCreateAccount: '%s' \n", err.Error()))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		jsonEncoder := json.NewEncoder(w)
		jsonEncoder.Encode(map[string]string{
			"token": token,
			"token_exp": exp,
			"refresh_token": refreshToken,
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

func (h *Handler) handleRefreshToken(w http.ResponseWriter, r *http.Request) {
	var err error

	// Parse the refresh token from the request
	var requestData struct {
		RefreshToken string `json:"refresh_token"`
		UserID string `json:"user_id"`
	}

	err = json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		log.Printf("ERROR in handleRefreshToken: %s \n", err)
		utils.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error in handleRefreshToken: '%s' \n", err.Error()))
		return
	}

	// Validate the refresh token
	claims, err := utils.ValidateRefreshToken(requestData.RefreshToken, requestData.UserID)
	if err != nil {
		log.Printf("ERROR in handleRefreshToken: %s \n", err)
		utils.RespondWithError(w, http.StatusUnauthorized, "Invalid or expired refresh token")
		return
	}

	// Generate a new JWT (access token) and possibly a new refresh token
	token, exp, err := utils.GenerateJWTToken(claims.UserID)
	if err != nil {
		log.Printf("ERROR in handleRefreshToken: %s \n", err)
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error in handleRefreshToken: '%s' \n", err.Error()))
		return
	}

	// // Optionally generate a new refresh token if desired
	// refreshToken, err := utils.GenerateRefreshToken(claims.UserID)
	// if err != nil {
	// 	log.Printf("ERROR in handleRefreshToken: %s \n", err)
	// 	utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error in handleRefreshToken: '%s' \n", err.Error()))
	// 	return
	// }

	// Send the new tokens back in the response
	w.Header().Set("Content-Type", "application/json")
	jsonEncoder := json.NewEncoder(w)
	jsonEncoder.Encode(map[string]string{
		"token":        token,
		"token_exp":    exp,
		// "refresh_token": refreshToken,
		"uuid":         claims.UserID,
	})
}
