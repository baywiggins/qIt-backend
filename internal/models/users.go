package models

import (
	"database/sql"
	"fmt"
)

type User struct {
	ID string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	UserState string `json:"state"`
}

func GetByUserName(db *sql.DB, username string) (User, error) {
	var err error;
	var user User;

	row := db.QueryRow(fmt.Sprintf("SELECT id, username, pass, user_state FROM Users WHERE username = '%s'", username))

	
	switch err = row.Scan(&user.ID, &user.Username, &user.Password, &user.UserState); err {
		case sql.ErrNoRows:
			return User{}, fmt.Errorf("'%s' does not exist", username)
		case nil:
			return user, err
		default:
			return User{}, err
	}
}

func InsertUser(db *sql.DB, user User) (error) {
	var err error;
	_, err = db.Exec("INSERT INTO Users (id, username, pass, user_state) VALUES(?, ?, ?, ?);", user.ID, user.Username, user.Password, user.UserState)
	return err
}