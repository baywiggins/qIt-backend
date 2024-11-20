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
			log.Printf("ERROR: %s \n", errors.New("unauthorized"))
			utils.RespondWithError(w, http.StatusUnauthorized, fmt.Sprintf("Error in AuthMiddleware: '%s'", errors.New("unauthorized")))
			return
		}
		_, err := utils.ValidateJWTToken(tokenString)
		if err != nil {
			log.Printf("ERROR: %s \n", err)
				utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error in AuthMiddleware: '%s' \n", err.Error()))
				return
		}

		next.ServeHTTP(w, r)
	})
}