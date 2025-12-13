package migration

import (
	"log"

	sqlengine "cozeva.com/vault/sql"
)

func InitSqliteDB() {
	db, err := sqlengine.GetSqlInstance()
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	initDbScheme := `
	CREATE TABLE IF NOT EXISTS project ( 
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE,
		uid TEXT UNIQUE
	);

	CREATE TABLE IF NOT EXISTS secret (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS secret_project_map (
		project_id INTEGER NOT NULL,
		secret_id INTEGER NOT NULL,
		PRIMARY KEY (project_id, secret_id),
		FOREIGN KEY (project_id) REFERENCES project(id)
			ON DELETE CASCADE ON UPDATE CASCADE,
		FOREIGN KEY (secret_id) REFERENCES secret(id)
			ON DELETE CASCADE ON UPDATE CASCADE
	);

	CREATE TABLE IF NOT EXISTS secret_version (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		secret_id INTEGER NOT NULL,
		environment TEXT NOT NULL,
		version INTEGER NOT NULL,
		value TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(secret_id, environment, version),
		FOREIGN KEY(secret_id) REFERENCES secret(id)
	);

	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
        first_name TEXT NOT NULL,
        last_name TEXT NOT NULL,
        email TEXT NOT NULL UNIQUE,
        password_hash TEXT NOT NULL,
        profile_picture TEXT,
        created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	);
	
	CREATE TABLE IF NOT EXISTS user_sessions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		session_id TEXT NOT NULL,
		refresh_token_hash TEXT NOT NULL,
        refresh_token_expires_at DATETIME NOT NULL,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id)
			ON DELETE CASCADE ON UPDATE CASCADE
	);
	`

	_, err = db.Exec(initDbScheme)
	if err != nil {
		log.Fatalf("Failed to create tables: %v", err)
	}
}
