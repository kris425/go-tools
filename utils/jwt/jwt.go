package jwt

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type Claims struct {
	Payload interface{} `json:"payload"`
	jwt.StandardClaims
}

// GenerateToken generate tokens used for auth
func GenerateToken(secret []byte, payload interface{}, expire time.Duration) (string, error) {
	expireTime := time.Now().Add(expire)

	claims := Claims{
		payload,
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "api",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(secret)

	return token, err
}

// ParseToken parsing token
func ParseToken(token string, secret []byte) (interface{}, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		return nil, err
	}
	if tokenClaims == nil || !tokenClaims.Valid {
		return nil, errors.New("Invalid token")
	}
	if claims, ok := tokenClaims.Claims.(*Claims); ok {
		return claims.Payload, nil
	}

	return nil, errors.New("Invalid token")
}
