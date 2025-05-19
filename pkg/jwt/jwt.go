package jwt

import (
	"fmt"
	"github.com/VladimirGladky/Yandex-Calculator/internal/domain/models"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

func NewToken(user *models.User, duration time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
		"nbf":   time.Now().Unix(),
		"exp":   time.Now().Add(duration).Unix(),
		"iat":   time.Now().Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		return "", fmt.Errorf("error creating token: %v", err)
	}
	return tokenString, nil
}

func ParseToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET")), nil
	})
}

func GetEmailFromToken(tokenString string) (string, error) {
	token, err := ParseToken(tokenString)
	if err != nil {
		return "", fmt.Errorf("error parsing token: %v", err)
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		email, ok := claims["email"].(string)
		if !ok {
			return "", fmt.Errorf("invalid token")
		}
		return email, nil
	}
	return "", fmt.Errorf("invalid token")
}
