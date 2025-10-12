package sqlengine

import (
	"database/sql"
	"os"

	"shounak.me/configmanager/config"
)

func GetSqlInstance() (*sql.DB, error) {
	dbDriver := os.Getenv("api.dbdriver")
	if dbDriver == "" {
		dbDriver = config.DefaultDBDriver
	}

	var db *sql.DB
	var err error

	switch dbDriver {
	case "sqlite":
		db, err = sql.Open("sqlite3", config.SqliteDbLocation)
	case "mysql":
	}

	if err != nil {
		return nil, err
	}

	return db, nil
}
