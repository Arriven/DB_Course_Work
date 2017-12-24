CREATE DATABASE course_db;
\c course_db;
CREATE TABLE IF NOT EXISTS phonebook("id" SERIAL PRIMARY KEY, "name" varchar(50), "phone" varchar(100));
DELETE FROM phonebook;
INSERT INTO phonebook VALUES (default, 'SomeName', '0123456789');
CREATE TABLE IF NOT EXISTS 
users(
    user_id SERIAL PRIMARY KEY,
    nickname varchar(20)
);
DELETE FROM users;
CREATE TABLE IF NOT EXISTS
projects(
    project_id SERIAL PRIMARY KEY,
    project_name varchar(50),
    owner SERIAL REFERENCES users(user_id)
);
DELETE FROM projects;
CREATE TYPE pull_request_status AS ENUM('pending', 'rejected', 'approved');
CREATE TABLE IF NOT EXISTS
pull_requests(
    id SERIAL PRIMARY KEY,
    project_id SERIAL REFERENCES projects(project_id),
    user_id SERIAL REFERENCES users(user_id),
    message varchar(500),
    status pull_request_status default 'pending'
);
DELETE FROM pull_requests;