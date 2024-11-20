package models

import (
	"database/sql"
	"time"
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