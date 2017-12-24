package main

import (
	"errors"
)

type Branch struct {
	id int64
	name string
	project Project
}

func (project Project)CreateBranch(name string) (*Branch, error) {
	var branch Branch
	err := project.owner.server.database.QueryRow("INSERT INTO branches VALUES(default, $1, $2) RETURNING branch_id", project.id, name).Scan(&branch.id)
	if err != nil {
		return nil, err
	}
	branch.project = project
	branch.name = name
	return &branch, err
}

func (server *Server)GetBranchById(id int64) (*Branch, error) {
	var branch Branch
	rows, err := server.database.Query("SELECT * FROM branches WHERE branch_id=$1", id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var project_id int64
		err = rows.Scan(&branch.id, &project_id, &branch.name)
		if err != nil {
			return nil, err
		}
		project, err := server.GetProjectById(project_id)
		if err != nil {
			return nil, err
		}
		branch.project = *project
		return &branch, err
	}
	return nil, errors.New("No branches found")
}

func (project Project) GetBranches() ([]Branch, error) {
	rows, err := project.owner.server.database.Query("SELECT * FROM branches WHERE branch_project=$1", project.id)
	if err != nil {
		return nil, err
	}
	var result []Branch

	for rows.Next() {
		var project_id int64
		var branch Branch
		err = rows.Scan(&branch.id, &project_id, &branch.name)
		if err != nil {
			return nil, err
		}
		if project.id != project_id {
			return nil, errors.New("Wrong project")
		}
		branch.project = project
		result = append(result, branch)
	}
	return result, nil
}

func (project Project) GetBranchByName(name string) (*Branch, error) {
	var branch Branch
	rows, err := project.owner.server.database.Query("SELECT * FROM branches WHERE branch_project=$1 AND branch_name=$2", project.id, name)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var project_id int64
		err = rows.Scan(&branch.id, &project_id, &branch.name)
		if err != nil {
			return nil, err
		}
		if project.id != project_id {
			return nil, errors.New("Wrong project")
		}
		branch.project = project
		return &branch, nil
	}
	return nil, errors.New("No Branches Found")
}