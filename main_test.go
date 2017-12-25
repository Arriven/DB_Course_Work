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

func TestNormalFlow(t *testing.T) {
	server, err := CreateServer("postgres", "", "course_db")
	if err != nil {
		t.Error(err)
		return
	}
	defer server.Shutdown()
	
	user, err := server.CreateUser("TestUser")
	if err != nil {
		t.Error(err)
		return
	}
	project, err := user.CreateProject("TestProject")
	if err != nil {
		t.Error(err)
		return
	}
	test, err := project.AddTest("/bin/true", "always pass")
	if err != nil {
		t.Error(err)
		return
	}
	branch, err := project.GetBranchByName("master")
	if err != nil {
		t.Error(err)
		return
	}
	commit, err := user.MakeCommit(*branch, "Test commit")
	if err != nil {
		t.Error(err)
		return
	}
	_, err = commit.RunTest(*test)
	if err != nil {
		t.Error(err)
		return
	}
	pr, err := server.CreatePullRequest(*commit, "Test pull request")
	if err != nil {
		t.Error(err)
		return
	}
	err = pr.Validate()
	if err != nil {
		t.Error(err)
		return
	}
	if pr.status != approved {
		t.Error("Wrong pr status")
		return
	}
}

func TestServer(t *testing.T) {
	server, err := CreateServer("postgres", "", "course_db")
	if err != nil {
		t.Error(err)
		return
	}
	defer server.Shutdown()
	_, err = CreateServer("postgres", "wrong password", "wrong_db")
	if err == nil {
		t.Error("Created server with invalid credentials")
		return
	}
}

func TestUser(t *testing.T) {
	server, err := CreateServer("postgres", "", "course_db")
	if err != nil {
		t.Error(err)
		return
	}
	defer server.Shutdown()

	_, err = server.GetUserByNickname("Arriven")
	if err == nil {
		t.Error("Found user when there's should be no user")
		return
	}
	_, err = server.GetUserById(582341024)
	if err == nil {
		t.Error("Found user when there's should be no user")
		return
	}
	
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
	_, err = server.CreateUser("Arriven")
	if err == nil {
		t.Error("Created user with same name")
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
	project, err = server.GetProjectById(project.id)
	if err != nil {
		t.Error(err)
		return
	}
	_, err = server.GetProjectById(-1)
	if err == nil {
		t.Error("Found project where shouldn't")
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

func TestCommits(t *testing.T) {
	server, err := CreateServer("postgres", "", "course_db")
	if err != nil {
		t.Error(err)
		return
	}
	defer server.Shutdown()
	
	user, err := server.CreateUser("SomeUser3")
	if err != nil {
		t.Error(err)
		return
	}
	project, err := user.CreateProject("testProject")
	if err != nil {
		t.Error(err)
		return
	}
	branch, err := project.GetBranchByName("master")
	if err != nil {
		t.Error(err)
		return
	}
	_, err = user.MakeCommit(*branch, "Test commit")
	if err != nil {
		t.Error(err)
		return
	}
}

func TestTests(t *testing.T) {
	server, err := CreateServer("postgres", "", "course_db")
	if err != nil {
		t.Error(err)
		return
	}
	defer server.Shutdown()
	
	user, err := server.CreateUser("SomeUser4")
	if err != nil {
		t.Error(err)
		return
	}
	project, err := user.CreateProject("testProject")
	if err != nil {
		t.Error(err)
		return
	}
	test, err := project.AddTest("/bin/true", "always pass")
	if err != nil {
		t.Error(err)
		return
	}
	tests, err := project.GetTests()
	if err != nil {
		t.Error(err)
		return
	}
	if len(tests) != 1 {
		t.Error("Wrong number of tests")
		return
	}
	branch, err := project.GetBranchByName("master")
	if err != nil {
		t.Error(err)
		return
	}
	commit, err := user.MakeCommit(*branch, "Test commit")
	if err != nil {
		t.Error(err)
		return
	}
	_, err = commit.RunTest(*test)
	if err != nil {
		t.Error(err)
		return
	}
	success, err := commit.RunTest(*test)
	if err != nil {
		t.Error(err)
		return
	}
	if !success {
		t.Error("Test failed")
		return
	}
	all_tests_success, err := commit.IsAllTestsPassed()
	if err != nil {
		t.Error(err)
		return
	}
	if !all_tests_success {
		t.Error("Some of the tests failed")
		return
	}
}