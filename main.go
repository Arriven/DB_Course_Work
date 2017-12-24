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

func clearDB(db *sql.DB) {
	fmt.Println("# Clearing Table")
	_, err := db.Exec("DELETE FROM phonebook")
	if(err != nil){
		fmt.Println(err)
		return
	}
}

func initDB(db *sql.DB) {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS " +
		"phonebook(\"id\" SERIAL PRIMARY KEY, " +
		"\"name\" varchar(50), \"phone\" varchar(100))")
	if(err != nil){
		fmt.Println(err)
		return
	}
	clearDB(db)
	fmt.Println("# Inserting Values")
	stmt, err := db.Prepare("INSERT INTO phonebook VALUES (default, $1, $2)")
	_, err = stmt.Exec("Bohdan", "0935293993")
	if(err != nil){
		fmt.Println(err)
		return
	}
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
	initDB(db)
}
