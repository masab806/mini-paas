package config

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	ID    int    `json:"user_id"`
	Email string `json:"email"`
	Username string `json:"username"`

	jwt.RegisteredClaims
}