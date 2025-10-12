package security

import (
	"strings"

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
