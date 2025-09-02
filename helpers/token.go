package helpers

import (
"fmt"
"os"
"time"

"github.com/golang-jwt/jwt/v5"
)

func GenerateResetToken(email string) (string, error) {
	secret := os.Getenv("RESET_PASSWORD_SECRET")
	expiry, _ := time.ParseDuration(os.Getenv("RESET_PASSWORD_EXPIRY"))

	claims := jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(expiry).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ParseResetToken(tokenString string) (string, error) {
	secret := os.Getenv("RESET_PASSWORD_SECRET")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(secret), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		email := claims["email"].(string)
		return email, nil
	}
	return "", err
}
