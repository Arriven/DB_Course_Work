package main

import (
	"errors"
)

type Project struct {
	id int64
	owner User
	name string
}

func (owner User) CreateProject (name string) (*Project, error) {
	var project Project
	err := owner.server.database.QueryRow("INSERT INTO projects(project_id, project_owner, project_name)" +
	" VALUES(default, $1, $2) RETURNING project_id",
	owner.id, name).Scan(&project.id)
	if err != nil {
		return nil, err
	}
	project.owner = owner
	project.name = name
	project.CreateBranch("master")
	return &project, err
}

func (server *Server) GetProjectById (id int64) (*Project, error) {
	var project Project
	rows, err := server.database.Query("SELECT project_id, project_owner, project_name" +
		" FROM projects WHERE project_id=$1", id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var owner_id int64
		err = rows.Scan(&project.id, &owner_id, &project.name)
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

func (user User) GetProjects() ([]Project, error) {
	rows, err := user.server.database.Query("SELECT project_id, project_owner, project_name" +
		" FROM projects WHERE project_owner=$1", user.id)
	if err != nil {
		return nil, err
	}
	var result []Project

	for rows.Next() {
		var owner_id int64
		var project Project
		err = rows.Scan(&project.id, &owner_id, &project.name)
		if err != nil {
			return nil, err
		}
		if user.id != owner_id {
			return nil, errors.New("Wrong owner")
		}
		project.owner = user
		result = append(result, project)
	}
	return result, nil
}