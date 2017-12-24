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
		return
	}
	defer server.Shutdown()
	
	user, err := server.CreateUser("Arriven")
	if err != nil {
		t.Error(err)
		return
	}
	if user.nickname != "Arriven" {
		t.Error("Warning: Name corrupted")
		return
	}
	user, err = server.GetUserByNickname("Arriven")
	if err != nil {
		t.Error(err)
		return
	}
	if user.nickname != "Arriven" {
		t.Error("Warning: Name corrupted")
		return
	}
	user, err = server.GetUserById(user.id)
	if err != nil {
		t.Error(err)
		return
	}
	if user.nickname != "Arriven" {
		t.Error("Warning: Name corrupted")
		return
	}
}

func TestProject(t *testing.T) {
	server, err := CreateServer("postgres", "", "course_db")
	if err != nil {
		t.Error(err)
		return
	}
	defer server.Shutdown()
	
	user, err := server.CreateUser("SomeUser")
	if err != nil {
		t.Error(err)
		return
	}
	project, err := user.CreateProject("testProject")
	if err != nil {
		t.Error(err)
		return
	}
	if project.name != "testProject" {
		t.Error("Warning: Name corrupted")
		return
	}
	if project.owner.id != user.id {
		t.Error("Warning: Owner corrupted")
		return
	}
	projects, err := user.GetProjects()
	if err != nil {
		t.Error(err)
		return
	}
	if len(projects) != 1 {
		t.Error("Wrong number of projects")
		return
	}
}

func TestBranches(t *testing.T) {
	server, err := CreateServer("postgres", "", "course_db")
	if err != nil {
		t.Error(err)
		return
	}
	defer server.Shutdown()
	
	user, err := server.CreateUser("SomeUser2")
	if err != nil {
		t.Error(err)
		return
	}
	project, err := user.CreateProject("testProject")
	if err != nil {
		t.Error(err)
		return
	}
	branches, err := project.GetBranches()
	if err != nil {
		t.Error(err)
		return
	}
	if len(branches) != 1 {
		t.Error("Wrong number of branches")
		return
	}
	branch, err := project.GetBranchByName("master")
	if err != nil {
		t.Error(err)
		return
	}
	if branch == nil {
		t.Error("master branch wasn't created")
		return
	}
}