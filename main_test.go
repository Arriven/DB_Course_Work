package main

import (
	"testing"
	"os"
	"fmt"
	"io/ioutil"
	"strings"
)

func checkErrorAndExit(err error) {
	if(err != nil) {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
}

func checkErrorAndWarning(err error) {
	if err != nil {
		fmt.Println("Warning: ", err)
	}
}

func TestMain(m *testing.M) {
	file, err := ioutil.ReadFile("dbsetup.sql")
	checkErrorAndExit(err)
	db, err := OpenDb("postgres", "", "")
	checkErrorAndExit(err)
	_, err = db.Exec("CREATE DATABASE course_db")
	checkErrorAndWarning(err)
	db, err = OpenDb("postgres", "", "course_db")
	checkErrorAndExit(err)
	requests := strings.Split(string(file), ";")
	for _, request := range requests {
		_, err := db.Exec(request)
		checkErrorAndWarning(err)
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
