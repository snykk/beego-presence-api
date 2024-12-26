package helpers

import (
	"errors"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtSecret = []byte("my-secret-tralalala")

// GenerateJWT generates a JWT token for a user with the given role
func GenerateJWT(userId int, email, role string) (string, error) {
	claims := jwt.MapClaims{
		"sub":   userId,
		"email": email,
		"role":  role,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
		"iat":   time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(jwtSecret)
	if err != nil {
		log.Println("Error signing token:", err)
		return "", err
	}
	return signedToken, nil
}

// ParseJWT parses the JWT token and returns the claims
func ParseJWT(tokenString string) (claims jwt.MapClaims, err error) {
	if token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	}); err != nil || !token.Valid {
		return nil, errors.New("token is not valid")
	}

	return
}

func GetRoleFromToken(tokenString string) (string, error) {
	claims, err := ParseJWT(tokenString)
	if err != nil {
		return "", err
	}

	if role, exists := claims["role"]; exists {
		return role.(string), nil
	}

	return "", errors.New("role not found in token")
}
