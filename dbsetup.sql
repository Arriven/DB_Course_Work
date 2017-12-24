CREATE TABLE IF NOT EXISTS phonebook("id" SERIAL PRIMARY KEY, "name" VARCHAR(50), "phone" VARCHAR(100));
DELETE FROM phonebook;
INSERT INTO phonebook VALUES (default, 'SomeName', '0123456789');
CREATE TABLE IF NOT EXISTS 
users(
    user_id SERIAL PRIMARY KEY,
    user_nickname VARCHAR(20) NOT NULL UNIQUE
);
DELETE FROM users;
CREATE TABLE IF NOT EXISTS
projects(
    project_id SERIAL PRIMARY KEY,
    project_name VARCHAR(50) NOT NULL,
    project_owner SERIAL REFERENCES users(user_id) NOT NULL
);
DELETE FROM projects;
CREATE TABLE IF NOT EXISTS
branches(
    branch_project SERIAL REFERENCES projects(project_id) NOT NULL,
    branch_name VARCHAR(10) NOT NULL,
    PRIMARY KEY (branch_project, branch_name)
);
DELETE FROM branches;
CREATE TABLE IF NOT EXISTS
commits(
    commit_id SERIAL PRIMARY KEY,
    commit_branch SERIAL REFERENCES branches(branch_id) NOT NULL,
    commit_author SERIAL REFERENCES users(user_id) NOT NULL,
    commit_message VARCHAR (500)
);
DELETE FROM commits;
CREATE TYPE pr_status AS ENUM('pending', 'rejected', 'approved');
CREATE TABLE IF NOT EXISTS
pull_requests(
    pull_request_id SERIAL PRIMARY KEY,
    pull_request_project SERIAL REFERENCES projects(project_id) NOT NULL,
    pull_request_commit SERIAL REFERENCES commits(commit_id) NOT NULL,
    pull_request_message VARCHAR(500),
    pull_request_status pr_status DEFAULT 'pending' NOT NULL
);
DELETE FROM pull_requests;
CREATE TABLE IF NOT EXISTS
tests(
    test_id SERIAL PRIMARY KEY,
    test_project SERIAL REFERENCES projects(project_id) NOT NULL,
    test_description VARCHAR(200)
);
DELETE FROM tests;
CREATE TABLE IF NOT EXISTS
test_results(
    test SERIAL REFERENCES tests(test_id) NOT NULL,
    commit SERIAL REFERENCES commits(commit_id) NOT NULL,
    success_status BOOLEAN DEFAULT FALSE NOT NULL,
    PRIMARY KEY(test, commit)
);
DELETE FROM test_results;