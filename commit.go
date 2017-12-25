package main

import (
	"errors"
)

type Commit struct {
	id int64
	author User
	branch Branch
	message string
}

func (user User) MakeCommit (branch Branch, message string) (*Commit, error) {
	var commit Commit
	err := user.server.database.QueryRow("INSERT INTO commits(commit_id, commit_author, commit_branch, commit_message)" +
		" VALUES(default, $1, $2, $3) RETURNING commit_id", user.id, branch.id, message).Scan(&commit.id)
	if err != nil {
		return nil, err
	}
	commit.author = user
	commit.branch = branch
	commit.message = message
	return &commit, nil
}

func (server *Server) GetCommitById(id int64) (*Commit, error) {
	var commit Commit
	rows, err := server.database.Query("SELECT commit_id, commit_author, commit_branch, commit_message" +
		" FROM commits WHERE commit_id=$1", id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var branch_id int64
		var author_id int64
		err = rows.Scan(&commit.id, &author_id, &branch_id, &commit.message)
		if err != nil {
			return nil, err
		}
		author, err := server.GetUserById(author_id)
		if err != nil {
			return nil, err
		}
		commit.author = *author
		branch, err := server.GetBranchById(branch_id)
		if err != nil {
			return nil, err
		}
		commit.branch = *branch
		return &commit, err
	}
	return nil, errors.New("No commits found")
}