package main

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type user struct {
	id int
	nickname string
}

func CreateUser(db *sql.DB, nickname string) (*user, error) {
	var newUser user
	err := db.QueryRow("INSERT INTO users VALUES(default, $1)", nickname).Scan(&newUser.id, &newUser.nickname)
	if err != nil {
		return nil, err
	}
	return &newUser, err
}