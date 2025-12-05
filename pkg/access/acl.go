package access

import (
	"log"

	"cozeva.com/vault/interfaces"
	sqlengine "cozeva.com/vault/sql"
)

func GetAccessForSecret(secretId any, userId string) (interfaces.AccessControlPayload, error) {
	var secretAccess interfaces.AccessControlPayload
	db, err := sqlengine.GetSqlInstance()
	if err != nil {
		return interfaces.AccessControlPayload{}, err
	}
	defer db.Close()

	rows, err := db.Query(sqlengine.CheckUserSecretAccessQuery, userId, secretId)

	if err != nil {
		return interfaces.AccessControlPayload{}, err
	}

	for rows.Next() {
		if err := rows.Scan(&secretAccess.HasReadAccess, &secretAccess.HasWriteAccess, &secretAccess.HasShareAccess, &secretAccess.HasDeleteAccess); err != nil {
			log.Fatal(err)
			return interfaces.AccessControlPayload{}, err
		}
	}

	return secretAccess, nil
}

func GetAccessForProject(projectId any, userId string) (interfaces.AccessControlPayload, error) {
	var projectAccess interfaces.AccessControlPayload
	db, err := sqlengine.GetSqlInstance()
	if err != nil {
		return interfaces.AccessControlPayload{}, err
	}
	defer db.Close()

	rows, err := db.Query(sqlengine.CheckUserProjectAccessQuery, userId, projectId)

	if err != nil {
		return interfaces.AccessControlPayload{}, err
	}

	for rows.Next() {
		if err := rows.Scan(&projectAccess.HasReadAccess, &projectAccess.HasWriteAccess, &projectAccess.HasShareAccess, &projectAccess.HasDeleteAccess); err != nil {
			log.Fatal(err)
			return interfaces.AccessControlPayload{}, err
		}
	}

	return projectAccess, nil
}
