package services

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

type AuthClaims struct {
	Address string `json:"address"`
	jwt.StandardClaims
}

func DecodeJwtToken(token string) (string, error) {
	var secretKey = []byte(os.Getenv("JWT_SECRET"))

	t, err := jwt.ParseWithClaims(token, &AuthClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected JWT signing method: %v", t.Header["alg"])
		}

		return secretKey, nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := t.Claims.(*AuthClaims); ok && t.Valid {
		return claims.Address, nil
	} else {
		return "", fmt.Errorf("invalid JWT token")
	}
}

func EncodeJwtToken(address string, expirationHour time.Duration) (string, error) {
	var secretKey = []byte(os.Getenv("JWT_SECRET"))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.MapClaims{
		"address": address,
		"exp":     time.Now().Add(expirationHour * time.Hour).Unix(),
	})

	return token.SignedString(secretKey)
}
