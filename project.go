package main

import (
	"errors"
)

type Project struct {
	id int64
	name string
	owner User
}

func (server *Server)CreateProject(name string, owner_id int64) (*Project, error) {
	var project Project
	err := server.database.QueryRow("INSERT INTO projects VALUES(default, $1, $2) RETURNING project_id", name, owner_id).Scan(&project.id)
	if err != nil {
		return nil, err
	}
	owner, err := server.GetUserById(owner_id)
	if err != nil {
		return nil, err
	}
	project.owner = *owner
	project.name = name
	return &project, err
}

func (server *Server)GetProjectById(id int64) (*Project, error) {
	var project Project
	rows, err := server.database.Query("SELECT * FROM projects WHERE project_id=$1", id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var owner_id int64
		err = rows.Scan(&project.id, &project.name, &owner_id)
		if err != nil {
			return nil, err
		}
		owner, err := server.GetUserById(owner_id)
		if err != nil {
			return nil, err
		}
		project.owner = *owner
		return &project, err
	}
	return nil, errors.New("No project found")
}