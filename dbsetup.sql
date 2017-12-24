DROP TABLE IF EXISTS test_results;
DROP TABLE IF EXISTS tests;
DROP TABLE IF EXISTS pull_requests;
DROP TYPE IF EXISTS pr_status;
DROP TABLE IF EXISTS commits;
DROP TABLE IF EXISTS branches;
DROP TABLE IF EXISTS projects;
DROP TABLE IF EXISTS users;
CREATE TABLE users(
    user_id SERIAL PRIMARY KEY,
    user_nickname VARCHAR(20) NOT NULL UNIQUE
);
CREATE TABLE projects(
    project_id SERIAL PRIMARY KEY,
    project_owner SERIAL REFERENCES users(user_id) NOT NULL,
    project_name VARCHAR(50) NOT NULL
);
CREATE TABLE branches(
    branch_id SERIAL PRIMARY KEY,
    branch_project SERIAL REFERENCES projects(project_id) NOT NULL,
    branch_name VARCHAR(10) NOT NULL,
    UNIQUE(branch_name, branch_project)
);
CREATE TABLE commits(
    commit_id SERIAL PRIMARY KEY,
    commit_branch SERIAL REFERENCES branches(branch_id) NOT NULL,
    commit_author SERIAL REFERENCES users(user_id) NOT NULL,
    commit_message VARCHAR (500)
);
CREATE TYPE pr_status AS ENUM('pending', 'rejected', 'approved');
CREATE TABLE pull_requests(
    pull_request_id SERIAL PRIMARY KEY,
    pull_request_project SERIAL REFERENCES projects(project_id) NOT NULL,
    pull_request_commit SERIAL REFERENCES commits(commit_id) NOT NULL,
    pull_request_message VARCHAR(500),
    pull_request_status pr_status DEFAULT 'pending' NOT NULL
);
CREATE TABLE tests(
    test_id SERIAL PRIMARY KEY,
    test_project SERIAL REFERENCES projects(project_id) NOT NULL,
    test_description VARCHAR(200)
);
CREATE TABLE test_results(
    test SERIAL REFERENCES tests(test_id) NOT NULL,
    commit SERIAL REFERENCES commits(commit_id) NOT NULL,
    success_status BOOLEAN DEFAULT FALSE NOT NULL,
    PRIMARY KEY(test, commit)
);