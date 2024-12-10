package usecase

import (
	"encoding/base64"
	"errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"strconv"
)

func HashToken(token string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(token), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func CompareTokenAndHash(token string, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(token))
}

func EncodeToken(token string) string {
	return base64.StdEncoding.EncodeToString([]byte(token))
}

func DecodeToken(token string) (string, error) {
	result, err := base64.StdEncoding.DecodeString(token)
	return string(result), err
}

func ValidateUserID(userID string) error {
	_, err := uuid.Parse(userID)
	return err
}

func StrSecToInt(unixTime string) (int, error) {
	ut, err := strconv.Atoi(unixTime)
	if err != nil {
		return 0, errors.New("Parse token time failed")
	}
	return ut, nil
}
