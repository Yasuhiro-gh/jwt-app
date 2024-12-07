package auth

import (
	"errors"
	"fmt"
	"github.com/Yasuhiro-gh/jwt-app/internal/usecase"
	"github.com/golang-jwt/jwt/v5"
	"strconv"
	"time"
)

const SECRETKEY = "yasuhiro_gh"

type RefreshClaims struct {
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
	// Passing there constant IPAddr to `emulate` difference
	refreshToken := BuildRefreshToken("0.0.0.0")
	return accessToken, refreshToken, nil
}

func BuildAccessToken(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, &AccessClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 10)),
		},
		UserID:          userID,
		IPAddr:          "0.0.0.0",
		SomePrivacyInfo: time.Now().String(),
	})
	return token.SignedString([]byte(SECRETKEY))
}

func BuildRefreshToken(IPAddr string) string {
	now := time.Now().Unix()
	stringNow := strconv.FormatInt(now, 10)
	tokenWithPayload := stringNow + "\n" + IPAddr
	encodedToken := usecase.EncodeToken(tokenWithPayload)
	return encodedToken
}

func GetRefreshedTokens(claims *AccessClaims, IPAddr string) (string, string, error) {
	if claims.IPAddr != IPAddr {
		usecase.SendEmail()
	}

	return GenerateTokenPair(claims.UserID)
}

func GetAccessClaims(accessToken string) (*AccessClaims, error) {
	claims := &AccessClaims{}
	token, err := jwt.ParseWithClaims(accessToken, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(SECRETKEY), nil
	})

	if err != nil {
		return &AccessClaims{}, errors.New("access token parse error: " + err.Error())
	}
	if !token.Valid {
		return &AccessClaims{}, errors.New("invalid access token")
	}

	return claims, nil
}
