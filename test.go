package main

import (
	"errors"
)

type Test struct {
	id int64
	project Project
	script string
	description string
}

func (project Project) AddTest (script string, descr string) (*Test, error) {
	var test Test
	err := project.owner.server.database.QueryRow("INSERT INTO tests(test_id, test_project, test_script_path, test_description)" +
		" VALUES(default, $1, $2, $3) RETURNING test_id", project.id, script, descr).Scan(&test.id)
	if err != nil {
		return nil, err
	}
	test.project = project
	test.script = script
	test.description = descr
	return &test, err
}

func (server *Server) GetTestById (id int64) (*Test, error) {
	var test Test
	rows, err := server.database.Query("SELECT test_id, test_project, test_script_path, test_description" +
		" FROM tests WHERE test_id=$1", id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var project_id int64
		err = rows.Scan(&test.id, &project_id, &test.script, &test.description)
		if err != nil {
			return nil, err
		}
		project, err := server.GetProjectById(project_id)
		if err != nil {
			return nil, err
		}
		test.project = *project
		return &test, err
	}
	return nil, errors.New("No tests found")
}

func (project Project) GetTests() ([]Test, error) {
	rows, err := project.owner.server.database.Query("SELECT test_id, test_project, test_script_path, test_description" +
		" FROM tests WHERE test_project=$1", project.id)
	if err != nil {
		return nil, err
	}
	var result []Test

	for rows.Next() {
		var project_id int64
		var test Test
		err = rows.Scan(&test.id, &project_id, &test.script, &test.description)
		if err != nil {
			return nil, err
		}
		if project.id != project_id {
			return nil, errors.New("Wrong project")
		}
		test.project = project
		result = append(result, test)
	}
	return result, nil
}