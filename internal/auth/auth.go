package auth

import (
	"errors"
	"github.com/Yasuhiro-gh/jwt-app/internal/usecase"
	"github.com/golang-jwt/jwt/v5"
	"strconv"
	"strings"
	"time"
)

const SECRETKEY = "yasuhiro_gh"

type RefreshClaims struct {
	IPAddr    string
	ExpiresAt string
}

type AccessClaims struct {
	jwt.RegisteredClaims
	UserID          string
	IPAddr          string
	SomePrivacyInfo string
}

func GenerateTokenPair(userID string) (string, string, error) {
	// Passing there constant IPAddr to `emulate` difference
	refreshToken := BuildRefreshToken("0.0.0.0")
	accessToken, err := BuildAccessToken(userID, "0.0.0.0")
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

func BuildAccessToken(userID, IPAddr string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, &AccessClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 10)),
		},
		UserID:          userID,
		IPAddr:          IPAddr,
		SomePrivacyInfo: time.Now().String(),
	})
	return token.SignedString([]byte(SECRETKEY))
}

func BuildRefreshToken(IPAddr string) string {
	expireAt := time.Now().Add(time.Hour * 10).Unix()
	strExpiresAt := strconv.FormatInt(expireAt, 10)
	tokenWithPayload := strExpiresAt + "\n" + IPAddr
	return tokenWithPayload
}

func GetRefreshedTokens(claims *RefreshClaims, userID, IPAddr string) (string, string, error) {
	if claims.IPAddr != IPAddr {
		usecase.SendEmail()
	}

	return GenerateTokenPair(userID)
}

func GetRefreshClaims(refreshToken string) (*RefreshClaims, error) {
	splitedDecodedToken := strings.Split(refreshToken, "\n")
	if len(splitedDecodedToken) != 2 {
		return &RefreshClaims{}, errors.New("invalid refresh token")
	}

	claims := &RefreshClaims{ExpiresAt: splitedDecodedToken[0], IPAddr: splitedDecodedToken[1]}
	if claims.ExpiresAt == "" || claims.IPAddr == "" {
		return &RefreshClaims{}, errors.New("invalid refresh token: Empty payload")
	}

	return claims, nil
}
