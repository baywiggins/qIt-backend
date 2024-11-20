package models

import "database/sql"

type User struct {
	ID string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	UserState string `json:"state"`
}

func FetchAll(db *sql.DB) ([]User, error) {
	var err error;

	rows, err := db.Query("SELECT id, username, pass, user_state FROM Users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Username, &user.Password, &user.UserState); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, err
}

func InsertUser(db *sql.DB, user User) error {
	var err error;
	_, err = db.Exec("INSERT INTO Users (id, username, pass, user_state) VALUES(?, ?, ?, ?);", user.ID, user.Username, user.Password, user.UserState)
	return err
}