package main

import (
	"errors"
)

type User struct {
	id int64
	nickname string
	server *Server
}

func (server *Server)CreateUser(nickname string) (*User, error) {
	var user User
	err := server.database.QueryRow("INSERT INTO users VALUES(default, $1) RETURNING user_id, user_nickname", nickname).Scan(&user.id, &user.nickname)
	if err != nil {
		return nil, err
	}
	user.server = server
	return &user, err
}

func (server *Server)GetUserByNickname(nickname string) (*User, error) {
	var user User
	user.server = server
	rows, err := server.database.Query("SELECT * FROM users WHERE user_nickname=$1", nickname)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(&user.id, &user.nickname)
		if err != nil {
			return nil, err
		}
		return &user, err
	}
	return nil, errors.New("No user found")
}

func (server *Server)GetUserById(id int64) (*User, error) {
	var user User
	user.server = server
	rows, err := server.database.Query("SELECT * FROM users WHERE user_id=$1", id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(&user.id, &user.nickname)
		if err != nil {
			return nil, err
		}
		return &user, err
	}
	return nil, errors.New("No user found")
}