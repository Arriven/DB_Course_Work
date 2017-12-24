package main

import (
	"testing"
	"os"
	"os/exec"
)

func TestMain(m *testing.M){
	exec.Command("psql", "-f dbsetup.sql -U postgres").Run()
	exec.Command("sleep", "10").Run()
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
