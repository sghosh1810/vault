package config

import (
	"os"
)

func GetConfigValue(key string) string {
	data := os.Getenv(key)
	if data == "" {
		data = getDefaultValue(key)
	}
	return data
}

func getDefaultValue(key string) string {
	switch key {
	case "db.sqlitelocation":
		return "./secret.db"
	case "db.defaultdriver":
		return "sqlite"
	case "api.config.version":
		return "v1"
	case "api.config.port":
		return "8080"
	}
	return ""
}
