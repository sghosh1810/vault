package security

import (
	"os"
)

func GetMasterKeyRing() []byte {
	masterKey := os.Getenv("api.masterkey")
	if masterKey == "" {
		panic("Cannot read master key from environment variables. Ensure api.masterkey is set as an environment variable.")
	}
	return []byte(masterKey)
}

func GetJwtKeyRing() []byte {
	jwtKey := os.Getenv("api.jwtkey")
	if jwtKey == "" {
		panic("Cannot read jwt secret from environment variables. Ensure api.jwtkey is set as an environment variable.")
	}
	return []byte(jwtKey)

}
