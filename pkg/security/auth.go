package security

import (
	"strconv"
	"strings"
	"time"

	"cozeva.com/vault/interfaces"
	"github.com/golang-jwt/jwt/v5"
)

func VerifyJWT(authHeader string) (jwt.MapClaims, any) {

	if authHeader == "" {
		return nil, map[string]string{
			"status":  "fail",
			"message": "Authorization header is not set",
		}
	}

	authToken := strings.Split(authHeader, " ")
	if len(authToken) != 2 || authToken[0] != "Bearer" {
		return nil, map[string]string{
			"status":  "fail",
			"message": "Authorization token must be in format Bearer {token}",
		}
	}

	tokenString := authToken[1]
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(GetJwtKeyRing()), nil
	})

	if err != nil || !token.Valid {
		return nil, map[string]string{
			"status":  "fail",
			"message": "Invalid or expired token",
		}
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return nil, map[string]string{
			"status":  "fail",
			"message": "Failed to parse token.",
		}
	}

	return claims, nil
}

func GenerateJWT(userID int64, sessionID string, ttl time.Duration) (string, error) {
	secret := []byte(GetJwtKeyRing())
	now := time.Now()
	expiresAt := now.Add(ttl)

	claims := interfaces.JwtClaims{
		Uid:       userID,
		SessionID: sessionID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			Subject:   strconv.FormatInt(userID, 10),
			ID:        sessionID,
			Issuer:    "vault",
		},
	}

	// HS256 signed token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
