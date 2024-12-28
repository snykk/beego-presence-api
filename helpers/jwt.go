package helpers

import (
	"errors"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtIssuer = "beego-presence-api"
var jwtSecret = []byte("my-secret-tralalala")

// GenerateJWT generates a JWT token for a user with the given role
func GenerateJWT(userId int, email, role string) (string, error) {
	claims := jwt.MapClaims{
		"sub":   userId,
		"email": email,
		"role":  role,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
		"iat":   time.Now().Unix(),
		"iss":   jwtIssuer,
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

func GetRoleFromMapClaims(claims jwt.MapClaims) (string, error) {
	if role, exists := claims["role"]; exists {
		return role.(string), nil
	}

	return "", errors.New("role not found in claims")
}

func GetUseridFromMapClaims(claims jwt.MapClaims) (int, error) {
	if userId, exists := claims["sub"]; exists {
		return int(userId.(float64)), nil
	}

	return 0, errors.New("userId not found in claims")
}
