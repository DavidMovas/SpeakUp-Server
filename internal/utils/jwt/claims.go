package jwt

import "github.com/golang-jwt/jwt/v5"

type AccessClaims struct {
	jwt.RegisteredClaims
	UserID string `json:"user_id"`
}
