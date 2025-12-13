package security

import (
	"runtime"

	"github.com/alexedwards/argon2id"
)

// HashPassword generates an Argon2id hash for the given plain password.
func GenerateHash(secret string) (string, error) {
	argon2Config := &argon2id.Params{
		Memory:      128 * 1024, // 128 MB
		Iterations:  4,
		Parallelism: uint8(runtime.NumCPU()),
		SaltLength:  16,
		KeyLength:   32,
	}
	hash, err := argon2id.CreateHash(secret, argon2Config)
	if err != nil {
		return "", err
	}
	return hash, nil
}

// CheckPassword compares a plain password with a stored Argon2id hash.
func CompareHash(secret string, hash string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(secret, hash)
	if err != nil {
		return false, err
	}
	return match, nil
}
