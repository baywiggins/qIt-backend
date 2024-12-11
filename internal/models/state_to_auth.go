package models

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/baywiggins/qIt-backend/pkg/utils"
)

type StateToAuth struct {
	UserState string `json:"state"`
	AuthToken string `json:"auth_token"`
	RefreshToken string `json:"refresh_token"`
	ExpirationTime time.Time `json:"expiration_time"`
}

func InsertStateToAuth(db *sql.DB, sta StateToAuth) error {
	var err error;
	_, err = db.Exec("INSERT INTO State_to_auth (user_state, auth_token, refresh_token, expiration_time) VALUES(?, ?, ?, ?);", sta.UserState, sta.AuthToken, sta.RefreshToken, sta.ExpirationTime)
	return err
}

func GetAuthTokenByUser(db *sql.DB, username string) (string, error) {
	var err error;
	var token string;

	query := fmt.Sprintf(`
		SELECT s.auth_token FROM State_to_auth s
		JOIN Users u ON s.user_state = u.user_state
		WHERE u.username = '%s';
	`, username)

	row := db.QueryRow(query)
	if err = row.Scan(token); err != nil {
		return "", err
	}

	return token, err
}

func GetAuthTokenByID(db *sql.DB, uuid string) (string, error) {
	var err error;
	var token string;

	query := fmt.Sprintf(`
		SELECT s.auth_token FROM State_to_auth s
		JOIN Users u ON s.user_state = u.user_state
		WHERE u.id = '%s';
	`, uuid)

	row := db.QueryRow(query)
	if err = row.Scan(&token); err != nil {
		return "", err
	}

	// Decrypt token
	decryptedToken, err := utils.Decrypt(token)
	if err != nil {
		return "", err
	}

	// Return token + err
	return decryptedToken, err
}

// func GetRefreshTokenByUser(db *sql.DB, username string) (string, error) {
// 	var err error;

// 	query := fmt.Sprintf(`
// 		SELECT s.refresh_token FROM State_to_auth s
// 		JOIN Users u ON s.user_state = u.user_state
// 		WHERE u.username = '%s';
// 	`, username)

// }