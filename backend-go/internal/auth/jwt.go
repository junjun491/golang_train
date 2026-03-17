package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret []byte

const (
	issuer   = "otayori_api"
	audience = "otayori_frontend"
)

func InitJWT(secret string) {
	jwtSecret = []byte(secret)
}

func GenerateToken(role string, id int64) (string, error) {
	now := time.Now()
	exp := now.Add(24 * time.Hour)

	claims := jwt.MapClaims{
		"sub": fmt.Sprintf("%s:%d", role, id),
		"iat": now.Unix(),
		"exp": exp.Unix(),
		"iss": issuer,
		"aud": audience,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ParseToken(tokenString string) (jwt.MapClaims, error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, fmt.Errorf("unexpected signing method")
		}

		return jwtSecret, nil
	},
		jwt.WithAudience(audience),
		jwt.WithIssuer(issuer),
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}