package main

import (
	"fmt"
	"database/sql"
	_ "github.com/lib/pq"
)

type Server struct {
	database *sql.DB
}

func OpenDb(username string, password string, database string) (*sql.DB, error){
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		username, password, database)
	return sql.Open("postgres", dbinfo)
}

func CreateServer(username string, password string, database string) (*Server, error) {
	db, err := OpenDb(username, password, database)
	if err != nil {
		return nil, err
	}
	server := Server{db}
	return &server, err
}