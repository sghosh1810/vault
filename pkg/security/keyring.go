package security

import (
	"cozeva.com/vault/config"
)

func GetMasterKeyRing() []byte {
	masterKey := config.GetConfigValue("core.masterkey")
	if masterKey == "" {
		panic("Cannot read master key from environment variables. Ensure core.masterkey is set as an environment variable.")
	}
	return []byte(masterKey)
}

func GetJwtKeyRing() []byte {
	jwtKey := config.GetConfigValue("core.jwtkey")
	if jwtKey == "" {
		panic("Cannot read jwt secret from environment variables. Ensure core.jwtkey is set as an environment variable.")
	}
	return []byte(jwtKey)

}
