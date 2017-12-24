package main

import (
	"testing"
	"os"
	"fmt"
	"io/ioutil"
	"strings"
)

func TestMain(m *testing.M) {
	file, err := ioutil.ReadFile("dbsetup.sql")
	if(err != nil) {
		fmt.Println(err)
		os.Exit(1)
	}
	db, err := OpenDb("postgres", "", "")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	_, err = db.Exec("CREATE DATABASE course_db")
	if err != nil {
		fmt.Println(err)
	}
	db, err = OpenDb("postgres", "", "course_db")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	requests := strings.Split(string(file), ";")
	for _, request := range requests {
		_, err := db.Exec(request)
		if(err != nil) {
			fmt.Println(err)
		}
	}
	os.Exit(m.Run())
}

func TestNormalFlow(t *testing.T) {
	err := PrintDb("postgres", "", "course_db")
	if(err != nil) {
		t.Error(err)
	}
}

func TestError(t *testing.T) {
	err := PrintDb("postgres", "123", "wrong_name")
	if(err == nil) {
		t.Fail()
	}
}
