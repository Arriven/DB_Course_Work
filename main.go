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

func PrintTable(db *sql.DB) {
	rows, err := db.Query("SELECT * FROM phonebook")
	if (err != nil){
		fmt.Println(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var rec record
		err = rows.Scan(&rec.id, &rec.name, &rec.phone)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(rec)
	}
}

func main() {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		DB_USER, DB_PASSWORD, DB_NAME)
	db, err := sql.Open("postgres", dbinfo)
	if(err != nil){
		return;
	}
	defer db.Close()
	PrintTable(db)
}
