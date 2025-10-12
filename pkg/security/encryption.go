package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

// EncryptSecret encrypts a plaintext secret using AES-GCM.
// Returns a base64 encoded ciphertext.
func EncryptSecret(plainText string) (string, error) {
	block, err := aes.NewCipher(GetMasterKeyRing())
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Create a random nonce (unique for each encryption)
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// Seal appends the encrypted data + authentication tag
	cipherText := aesGCM.Seal(nonce, nonce, []byte(plainText), nil)
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

// DecryptSecret decrypts a base64 encoded AES-GCM ciphertext.
func DecryptSecret(encodedCipher string) (string, error) {
	cipherData, err := base64.StdEncoding.DecodeString(encodedCipher)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(GetMasterKeyRing())
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := aesGCM.NonceSize()
	if len(cipherData) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, cipherText := cipherData[:nonceSize], cipherData[nonceSize:]

	plainText, err := aesGCM.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return "", err
	}

	return string(plainText), nil
}
