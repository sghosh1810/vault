package helper

import (
	sqlquery "cozeva.com/vault/sql"
)

func CleanupExpiredSessions() error {
	db, err := sqlquery.GetSqlInstance()
	if err != nil {
		return err
	}
	defer db.Close()
	_, err = db.Exec(sqlquery.DeleteExpiredSessionsQuery)
	return err
}
