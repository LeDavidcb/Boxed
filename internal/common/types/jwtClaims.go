package types

import "github.com/golang-jwt/jwt/v5"

type ResponseClaims struct {
	Name string
	jwt.RegisteredClaims
}
