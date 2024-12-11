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
	FinishedCreating int `json:"finished_creating"`
}

func GetByUserName(db *sql.DB, username string) (User, error) {
	var err error;
	var user User;

	row := db.QueryRow(fmt.Sprintf("SELECT id, username, pass, user_state, finished_creating FROM Users WHERE username = '%s'", username))

	
	switch err = row.Scan(&user.ID, &user.Username, &user.Password, &user.UserState, &user.FinishedCreating); err {
		case sql.ErrNoRows:
			return User{}, fmt.Errorf("'%s' does not exist", username)
		case nil:
			return user, err
		default:
			return User{}, err
	}
}

func UpdateCreationStatusByState(db *sql.DB, state string) (error) {
	var err error;

	_, err = db.Exec("UPDATE Users SET finished_creating = ? WHERE user_state = ?", 1, state)

	return err
}

func InsertUser(db *sql.DB, user User) (error) {
	var err error;
	_, err = db.Exec("INSERT INTO Users (id, username, pass, user_state, finished_creating) VALUES(?, ?, ?, ?, ?);", user.ID, user.Username, user.Password, user.UserState, 0)
	return err
}

func DeleteUserByID(db *sql.DB, id string) (error) {
	var err error;

	_, err = db.Exec("DELETE FROM Users WHERE id = ?", id)

	return err
}