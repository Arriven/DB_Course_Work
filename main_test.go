package main

import (
	"testing"
	"os"
	"os/exec"
)

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
