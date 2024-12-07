package usecase

import (
	"encoding/base64"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func HashToken(token string) (string, error) {
	fmt.Println(token)
	hashed, err := bcrypt.GenerateFromPassword([]byte(token), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func CompareHashAndToken(token string, hashed string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(token))
}

func EncodeToken(token string) string {
	return base64.StdEncoding.EncodeToString([]byte(token))
}

func DecodeToken(token string) (string, error) {
	result, err := base64.StdEncoding.DecodeString(token)
	return string(result), err
}
