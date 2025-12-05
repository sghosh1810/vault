package sqlengine

import (
	"database/sql"

	"cozeva.com/vault/config"
)

func GetSqlInstance() (*sql.DB, error) {
	dbDriver := config.GetConfigValue("db.defaultdriver")

	var db *sql.DB
	var err error

	switch dbDriver {
	case "sqlite":
		db, err = sql.Open("sqlite3", config.GetConfigValue("db.sqlitelocation"))
	case "mysql":
	}

	if err != nil {
		return nil, err
	}

	return db, nil
}
