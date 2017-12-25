package main

import (
	"database/sql/driver"
	_ "github.com/lib/pq"
	"errors"
)

const (
	pending PullRequestStatus = "pending"
	rejected PullRequestStatus = "rejected"
	approved PullRequestStatus = "approved"
)

type PullRequestStatus string

func (s *PullRequestStatus) Scan(value interface{}) error {
	asBytes, ok := value.([]byte)
	if !ok {
		return errors.New("Scan source is not []byte")
	}
	*s = PullRequestStatus(string(asBytes))
	return nil
}

func (s PullRequestStatus) Value() (driver.Value, error) {
	// validation would go here
	return string(s), nil
}
type PullRequest struct {
	id int64
	commit Commit
	message string
	status PullRequestStatus
}

func (server *Server) CreatePullRequest (commit Commit, message string) (*PullRequest, error) {
	var pr PullRequest
	err := server.database.QueryRow("INSERT INTO pull_requests(pull_request_id, pull_request_commit, pull_request_message, pull_request_status)" +
	" VALUES(default, $1, $2, default) RETURNING pull_request_id, pull_request_status", commit.id, message).Scan(&pr.id, &pr.status)
	if err != nil {
		return nil, err
	}
	pr.commit = commit
	pr.message = message
	return &pr, err
}

func (pr *PullRequest) ValidatePullRequest() error {
	success, err := pr.commit.RunAllTests()
	if err != nil {
		return err
	}
	status := rejected
	if success {
		status = approved
	}
	err = pr.commit.branch.project.owner.server.database.QueryRow("UPDATE pull_requests" +
		"SET pull_request_status=$2 WHERE pull_request_id=$1 RETURNING pull_request_status", pr.id, status).Scan(&pr.status)
	if err != nil {
		return err
	}
	return nil
}