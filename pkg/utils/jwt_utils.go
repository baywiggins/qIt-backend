package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var SecretKey = []byte("secret-key")
var RefreshSecretKey = []byte("refresh-secret-key") // Refresh token secret key

// Claims structure for both access and refresh tokens
type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

// GenerateJWTToken generates a short-lived access token (e.g., 15 minutes)
func GenerateJWTToken(userID string) (string, string, error) {
	expirationTime := jwt.NewNumericDate(time.Now().Add(1 * time.Hour)) // Access token expiration time
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "qIt",
			ExpiresAt: expirationTime,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(SecretKey)

	exp := time.Now().UTC().Add(59 * time.Minute).Format(time.RFC3339)

	return signed, exp, err
}

// GenerateRefreshToken generates a long-lived refresh token (e.g., 7 days)
func GenerateRefreshToken(userID string) (string, error) {
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "qIt",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(RefreshSecretKey)
}

// ValidateJWTToken validates the access token
func ValidateJWTToken(tokenString string, userID string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	if claims.UserID != userID {
		return nil, errors.New("unauthorized user")
	}

	// If token is expired, return error
	if claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, errors.New("access token expired")
	}

	return claims, nil
}

// ValidateRefreshToken validates the refresh token and returns the claims if valid
func ValidateRefreshToken(tokenString string, userID string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return RefreshSecretKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid refresh token")
	}
	if claims.UserID != userID {
		return nil, errors.New("unauthorized user")
	}

	// If the refresh token is expired, return error
	if claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, errors.New("refresh token expired")
	}

	return claims, nil
}

// RefreshAccessToken uses the refresh token to generate a new access token
func RefreshAccessToken(refreshTokenString string, userID string) (string, string, error) {
	// Validate the refresh token
	_, err := ValidateRefreshToken(refreshTokenString, userID)
	if err != nil {
		return "", "", err
	}

	// If the refresh token is valid, generate a new access token
	return GenerateJWTToken(userID)
}
