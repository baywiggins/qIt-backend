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
	ExpirationTime string `json:"expiration_time"`
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

func GetStateToAuthRowByID(db *sql.DB, uuid string) (string, string, time.Time, error) {
	var err error;
	var accessToken string;
	var refreshToken string;
	var expiresIn string;

	query := fmt.Sprintf(`
		SELECT s.auth_token, s.refresh_token, s.expiration_time FROM State_to_auth s
		JOIN Users u ON s.user_state = u.user_state
		WHERE u.id = '%s';
	`, uuid)

	row := db.QueryRow(query)
	if err = row.Scan(&accessToken, &refreshToken, &expiresIn); err != nil {
		return "", "", time.Now(), err
	}

	// Decrypt tokens
	decryptedAccessToken, err := utils.Decrypt(accessToken)
	if err != nil {
		return "", "", time.Now(), err
	}
	decryptedRefreshToken, err := utils.Decrypt(refreshToken)
	if err != nil {
		return "", "", time.Now(), err
	}

	expiresInTime, err := time.Parse(time.RFC3339, expiresIn)
	if err != nil {
		return "", "", time.Now(), err
	}

	return decryptedAccessToken, decryptedRefreshToken, expiresInTime, err
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

func UpdateAccessTokenByID(db *sql.DB, uuid string, accessToken string, refreshToken string, expirationTime string) (error) {
	var err error;

	// Encrypt our tokens
	encryptedAuthToken, err := utils.Encrypt(accessToken)
	if err != nil {
		return err
	}
	encryptedRefreshToken, err := utils.Encrypt(refreshToken)
	if err != nil {
		return err
	}
	
	// Update access token + refresh token
	query := fmt.Sprintf(`
		UPDATE State_to_auth
		SET
			auth_token = '%s',
			refresh_token = '%s',
			expiration_time = '%s'
		WHERE user_state IN (
			SELECT user_state FROM Users
			WHERE id = '%s'
		);
	`, encryptedAuthToken, encryptedRefreshToken, expirationTime, uuid)
	// Executre query and return error if exists
	_, err = db.Exec(query)

	fmt.Println(query)
	fmt.Println(err)

	return err
}

// func GetRefreshTokenByUser(db *sql.DB, username string) (string, error) {
// 	var err error;

// 	query := fmt.Sprintf(`
// 		SELECT s.refresh_token FROM State_to_auth s
// 		JOIN Users u ON s.user_state = u.user_state
// 		WHERE u.username = '%s';
// 	`, username)

// }