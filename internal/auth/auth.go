package auth

import (
	"errors"
	"fmt"
	"github.com/Yasuhiro-gh/jwt-app/internal/usecase"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"time"
)

const SECRETKEY = "yasuhiro_gh"

type RefreshClaims struct {
	jwt.RegisteredClaims
	IPAddr string
	UserID string
}

type AccessClaims struct {
	jwt.RegisteredClaims
	UserID          string
	IPAddr          string
	SomePrivacyInfo string
}

func GenerateTokenPair(userID string) (string, string, error) {
	accessToken, err := BuildAccessToken(userID)
	if err != nil {
		return "", "", err
	}
	refreshToken, err := BuildRefreshToken(userID)
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

func BuildAccessToken(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, &AccessClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 10)),
		},
		UserID:          userID,
		IPAddr:          "0.0.0.0:0000",
		SomePrivacyInfo: time.Now().String(),
	})
	return token.SignedString([]byte(SECRETKEY))
}

func BuildRefreshToken(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, &RefreshClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 2)),
		},
		UserID: userID,
		IPAddr: "0.0.0.0",
	})
	return token.SignedString([]byte(SECRETKEY))
}

func GetClaimsFromRefreshToken(refreshToken string) (*RefreshClaims, error) {
	claims := &RefreshClaims{}
	token, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(SECRETKEY), nil
	})

	if err != nil {
		return &RefreshClaims{}, errors.New("token parse error")
	}

	if !token.Valid || claims.UserID == "" || claims.IPAddr == "" {
		return &RefreshClaims{}, errors.New("token is invalid")
	}

	return claims, nil
}

func GetRefreshedTokens(claims *RefreshClaims, IPAddr string) (string, string, error) {
	if claims.IPAddr != IPAddr {
		usecase.SendEmail()
	}

	return GenerateTokenPair(claims.UserID)
}

func HashToken(token string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(token), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func CompareHashAndToken(token string, hashed string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(token))
}
