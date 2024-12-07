package auth

import (
	"errors"
	"fmt"
	"github.com/Yasuhiro-gh/jwt-app/internal/usecase"
	"github.com/golang-jwt/jwt/v5"
	"strings"
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
	refreshToken := BuildRefreshToken(userID)
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

func BuildRefreshToken(userID string) string {
	IPAddr := "0.0.0.0"
	tokenWithPayload := userID + "\n" + IPAddr
	encodedToken := usecase.EncodeToken(tokenWithPayload)
	return encodedToken
}

func GetClaimsFromRefreshToken(refreshToken string) (*RefreshClaims, error) {
	decodedToken, err := usecase.DecodeToken(refreshToken)
	if err != nil {
		return &RefreshClaims{}, errors.New("token parse error: " + err.Error())
	}

	splitedDecodedToken := strings.Split(decodedToken, "\n")
	if len(splitedDecodedToken) != 2 {
		return &RefreshClaims{}, errors.New("invalid refresh token")
	}

	claims := &RefreshClaims{UserID: splitedDecodedToken[0], IPAddr: splitedDecodedToken[1]}
	if claims.UserID == "" || claims.IPAddr == "" {
		return &RefreshClaims{}, errors.New("invalid refresh token: Empty payload")
	}

	return claims, nil
}

func GetRefreshedTokens(claims *RefreshClaims, IPAddr string) (string, string, error) {
	if claims.IPAddr != IPAddr {
		usecase.SendEmail()
	}

	return GenerateTokenPair(claims.UserID)
}

func ValidateAccessToken(accessToken string) error {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(SECRETKEY), nil
	})

	if err != nil {
		return errors.New("token parse error: " + err.Error())
	}
	if !token.Valid {
		return errors.New("invalid access token")
	}

	return nil
}
