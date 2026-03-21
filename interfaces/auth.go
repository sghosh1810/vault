package interfaces

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtClaims struct {
	Uid       int64  `json:"uid"`
	SessionID string `json:"session_id"`
	jwt.RegisteredClaims
}

//	type UserRefreshTokenInfo struct {
//		ID                 int64     `json:"id"`
//		RefreshTokenHash   string    `json:"refresh_token"`
//		RefreshTokenExpiry time.Time `json:"refresh_token_expiry"`
//	}
type UserRefreshTokenInfo struct {
	ID                 int64     `json:"id"`
	UserID             int64     `json:"user_id"`
	RefreshTokenHash   string    `json:"refresh_token"`
	RefreshTokenExpiry time.Time `json:"refresh_token_expiry"`
}
