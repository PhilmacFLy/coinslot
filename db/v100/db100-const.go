package db100

const createSQLlitestmt = `
--
-- File generated with SQLiteStudio v3.0.7 on Do. Sep. 22 12:31:12 2016
--
-- Text encoding used: UTF-8
--
PRAGMA foreign_keys = off;
BEGIN TRANSACTION;

-- Table: Fundings
CREATE TABLE Fundings (user_id INTEGER NOT NULL, project_id INTEGER NOT NULL, amount DOUBLE NOT NULL, confirmed BOOLEAN NOT NULL DEFAULT (0), PRIMARY KEY (user_id, project_id) ON CONFLICT ROLLBACK);

-- Table: Users
CREATE TABLE Users (user_id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL, username VARCHAR (255) NOT NULL, password CHAR (64) NOT NULL, salt CHAR (64) NOT NULL, email VARCHAR (255) NOT NULL, "right" INTEGER NOT NULL);

-- Table: Projects
CREATE TABLE Projects (project_id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL, name STRING NOT NULL, goal DOUBLE NOT NULL, initiator INTEGER REFERENCES Users (user_id) ON DELETE CASCADE ON UPDATE CASCADE NOT NULL, description TEXT NOT NULL);

COMMIT TRANSACTION;
PRAGMA foreign_keys = on;
`
