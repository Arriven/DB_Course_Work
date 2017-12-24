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

func OpenDb(username string, password string, database string) (*sql.DB, error){
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		username, password, database)
	return sql.Open("postgres", dbinfo)
}

func PrintDb(username string, password string, database string) (error) {
	db, err := OpenDb(username, password, database)
	if(err != nil){
		return err;
	}
	defer db.Close()
	return PrintTable(db)
}

func main(){

}
