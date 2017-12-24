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
	defer db.Close()
	checkErrorAndExit(err)
	_, err = db.Exec("CREATE DATABASE course_db")
	checkErrorAndWarning(err)
	db, err = OpenDb("postgres", "", "course_db")
	checkErrorAndExit(err)
	defer db.Close()
	requests := strings.Split(string(file), ";")
	for _, request := range requests {
		_, err := db.Exec(request)
		checkErrorAndWarning(err)
	}
	os.Exit(m.Run())
}

func TestUser(t *testing.T) {
	server, err := CreateServer("postgres", "", "course_db")
	if err != nil {
		t.Error(err)
	}
	defer server.Shutdown()
	
	user, err := server.CreateUser("Arriven")
	if err != nil {
		t.Error(err)
	}
	if user.nickname != "Arriven" {
		t.Error("Warning: Name corrupted")
	}
	user, err = server.GetUserByNickname("Arriven")
	if err != nil {
		t.Error(err)
	}
	if user.nickname != "Arriven" {
		t.Error("Warning: Name corrupted")
	}
	user, err = server.GetUserById(user.id)
	if err != nil {
		t.Error(err)
	}
	if user.nickname != "Arriven" {
		t.Error("Warning: Name corrupted")
	}
}

func TestProject(t *testing.T) {
	server, err := CreateServer("postgres", "", "course_db")
	if err != nil {
		t.Error(err)
	}
	defer server.Shutdown()
	
	user, err := server.CreateUser("SomeUser")
	if err != nil {
		t.Error(err)
	}
	project, err := user.CreateProject("testProject")
	if err != nil {
		t.Error(err)
	}
	if project.name != "testProject" {
		t.Error("Warning: Name corrupted")
	}
	if project.owner.id != user.id {
		t.Error("Warning: Owner corrupted")
	}
	projects, err := user.GetProjects()
	if err != nil {
		t.Error(err)
	}
	if len(projects) != 1 {
		t.Error("Wrong number of projects")
	}
}