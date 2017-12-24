package main

import (
	"fmt"
	"database/sql"
	_ "github.com/lib/pq"
)

const (
	DB_USER 	= "postgres"
	DB_PASSWORD	= ""
	DB_NAME		= "course_db"
)

type record struct {
	id int
	name string
	phone string
}

func PrintTable(db *sql.DB) (error) {
	rows, err := db.Query("SELECT * FROM phonebook")
	if (err != nil){
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var rec record
		err = rows.Scan(&rec.id, &rec.name, &rec.phone)
		if err != nil {
			return err
		}
		fmt.Println(rec)
	}
	return nil
}

func PrintDb(username string, password string, database string) (error) {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		username, password, database)
	db, err := sql.Open("postgres", dbinfo)
	if(err != nil){
		return err;
	}
	defer db.Close()
	return PrintTable(db)
}
