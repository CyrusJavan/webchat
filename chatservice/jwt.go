package chatservice

import (
	"github.com/dgrijalva/jwt-go"
	"os"
)

func GetToken(m map[string]interface{}) (string, error) {
	key := []byte(os.Getenv("JWTKEY"))
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(m))

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(key)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}