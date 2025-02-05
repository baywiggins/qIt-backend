package middlewares

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/baywiggins/qIt-backend/pkg/utils"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			log.Printf("ERROR in AuthMiddleware: %s \n", errors.New("unauthorized"))
			utils.RespondWithError(w, http.StatusUnauthorized, fmt.Sprintf("Error: '%s'", errors.New("unauthorized")))
			return
		}

		userID := r.Header.Get("uuid")
		if userID == "" {
			log.Printf("ERROR in AuthMiddleware: %s \n", errors.New("unauthorized"))
			utils.RespondWithError(w, http.StatusUnauthorized, fmt.Sprintf("Error: '%s'", errors.New("unauthorized, missing uuid")))
			return
		}

		_, err := utils.ValidateJWTToken(tokenString, userID)
		if err != nil {
			statusCode := http.StatusUnauthorized
			if err.Error() != "invalid token" && err.Error() != "unauthorized user" && err.Error() != "token has invalid claims: token is expired" {
				statusCode = http.StatusInternalServerError
			}
			
			log.Printf("ERROR in AuthMiddleware: %s \n", err)
				utils.RespondWithError(w, statusCode, fmt.Sprintf("Error: '%s' \n", err.Error()))
				return
		}

		next.ServeHTTP(w, r)
	})
}