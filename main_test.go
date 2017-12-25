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
	} else {
		defer server.Shutdown()
	}
	_, err = CreateServer("postgres", "wrong password", "wrong_db")
	if err == nil {
		t.Error("Created server with invalid credentials")
	}
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
	} else if user.nickname != "Arriven" {
		t.Error("Warning: Name corrupted")
	}
	user, err = server.GetUserById(user.id)
	if err != nil {
		t.Error(err)
	} else if user.nickname != "Arriven" {
		t.Error("Warning: Name corrupted")
	}
	_, err = server.CreateUser("Arriven")
	if err == nil {
		t.Error("Created user with same name")
	}
	_, err = server.GetUserByNickname("WrongUserName")
	if err == nil {
		t.Error("Found user when there's should be no user")
	}
	_, err = server.GetUserById(-1)
	if err == nil {
		t.Error("Found user when there's should be no user")
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
	}
	if project.owner.id != user.id {
		t.Error("Warning: Owner corrupted")
	}
	projects, err := user.GetProjects()
	if err != nil {
		t.Error(err)
	} else if len(projects) != 1 {
		t.Error("Wrong number of projects")
	}
	project, err = server.GetProjectById(project.id)
	if err != nil {
		t.Error(err)
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
	} else if len(branches) != 1 {
		t.Error("Wrong number of branches")
	}
	branch, err := project.GetBranchByName("master")
	if err != nil {
		t.Error(err)
	}
	if branch == nil {
		t.Error("master branch wasn't created")
	}
	_, err = server.GetBranchById(branch.id)
	if err != nil {
		t.Error(err)
	}
	_, err = project.GetBranchByName("unknown branch")
	if err == nil {
		t.Error("Found branch when shouldn't")
	}
	_, err = server.GetBranchById(-1)
	if err == nil {
		t.Error("Found branch when shouldn't")
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
	commit, err := user.MakeCommit(*branch, "Test commit")
	if err != nil {
		t.Error(err)
		return
	}
	_, err = server.GetCommitById(commit.id)
	if err != nil {
		t.Error(err)
	}
	_, err = server.GetCommitById(-1)
	if err == nil {
		t.Error("Found commit when shouldn't")
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
	} else if len(tests) != 1 {
		t.Error("Wrong number of tests")
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
	success, err := commit.RunTest(*test)
	if err != nil {
		t.Error(err)
	}
	if !success {
		t.Error("Test failed")
	}
	all_tests_success, err := commit.IsAllTestsPassed()
	if err != nil {
		t.Error(err)
	}
	if !all_tests_success {
		t.Error("Some of the tests failed")
	}
	_, err = server.GetTestById(test.id)
	if err != nil {
		t.Error(err)
	}
	_, err = server.GetTestById(-1)
	if err == nil {
		t.Error("Found test where shouldn't")
	}
}

func TestPullRequest (t *testing.T) {
	server, err := CreateServer("postgres", "", "course_db")
	if err != nil {
		t.Error(err)
		return
	}
	defer server.Shutdown()
	
	user, err := server.CreateUser("SomeUser5")
	if err != nil {
		t.Error(err)
		return
	}
	project, err := user.CreateProject("testProject")
	if err != nil {
		t.Error(err)
		return
	}
	_, err = project.AddTest("/bin/true", "always pass")
	if err != nil {
		t.Error(err)
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
	pr, err := server.CreatePullRequest(*commit, "test pull request")
	if err != nil {
		t.Error(err)
		return
	}
	err = pr.Validate()
	if err != nil {
		t.Error(err)
	}
	_, err = server.GetPullRequestById(pr.id)
	if err != nil {
		t.Error(err)
	}
	_, err = server.GetPullRequestById(-1)
	if err == nil {
		t.Error("Found pull request when shouldn't")
	}
}