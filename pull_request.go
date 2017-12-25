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
	if s != pending && s != rejected && s != approved {
		return nil, errors.New("Invalid status")
	}
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

func (server *Server) GetPullRequestById(id int64) (*PullRequest, error) {
	var pr PullRequest
	rows, err := server.database.Query("SELECT pull_request_id, pull_request_commit, pull_request_message, pull_request_status" +
		" FROM pull_requests WHERE pull_request_id=$1", id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var commit_id int64
		err = rows.Scan(&pr.id, &commit_id, &pr.message, &pr.status)
		if err != nil {
			return nil, err
		}
		commit, err := server.GetCommitById(commit_id)
		if err != nil {
			return nil, err
		}
		pr.commit = *commit
		return &pr, err
	}
	return nil, errors.New("No pull requests found")
}

func (pr *PullRequest) Validate() error {
	success, err := pr.commit.RunAllTests()
	if err != nil {
		return err
	}
	status := rejected
	if success {
		status = approved
	}
	err = pr.commit.branch.project.owner.server.database.QueryRow("UPDATE pull_requests" +
		" SET pull_request_status = $1 WHERE pull_request_id = $2 RETURNING pull_request_status", status, pr.id).Scan(&pr.status)
	if err != nil {
		return err
	}
	return nil
}