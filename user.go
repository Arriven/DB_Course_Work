package main

import (
)

type User struct {
	id int64
	nickname string
}

func (server *Server)CreateUser(nickname string) (*User, error) {
	var newUser User
	err := server.database.QueryRow("INSERT INTO users VALUES(default, $1) RETURNING user_id, user_nickname", nickname).Scan(&newUser.id, &newUser.nickname)
	if err != nil {
		return nil, err
	}
	return &newUser, err
}