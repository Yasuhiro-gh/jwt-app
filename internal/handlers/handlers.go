package handlers

import (
	"encoding/json"
	"github.com/Yasuhiro-gh/jwt-app/internal/auth"
	"github.com/Yasuhiro-gh/jwt-app/internal/usecase"
	"net/http"
	"time"
)

type TokenHandler struct {
	usecase.TokenStore
}

func NewTokenHandler(ts *usecase.TokenStorage) *TokenHandler {
	return &TokenHandler{ts}
}

func Router(ts *usecase.TokenStorage) *http.ServeMux {
	mux := http.NewServeMux()
	th := NewTokenHandler(ts)

	mux.Handle("/api/tokens/{uid}", th.GenerateTokenPair())
	mux.Handle("/api/refresh", th.Refresh())

	return mux
}

func (th *TokenHandler) GenerateTokenPair() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Provide Get method", http.StatusMethodNotAllowed)
			return
		}

		userID := r.PathValue("uid")
		if userID == "" {
			http.Error(w, "Provide user id in parameters", http.StatusBadRequest)
			return
		}

		if err := usecase.ValidateUserID(userID); err != nil {
			http.Error(w, "Provide valid UUID user id in parameters", http.StatusBadRequest)
			return
		}

		accessToken, refreshToken, err := auth.GenerateTokenPair(userID)
		if err != nil {
			http.Error(w, "Can't generate jwt tokens: "+err.Error(), http.StatusInternalServerError)
			return
		}

		hashedToken, err := usecase.HashToken(refreshToken)
		if err != nil {
			http.Error(w, "Can't hash refresh token: "+err.Error(), http.StatusInternalServerError)
			return
		}
		err = th.SetNewToken(userID, hashedToken)
		if err != nil {
			http.Error(w, "Can't set new token: "+err.Error(), http.StatusInternalServerError)
			return
		}

		type tokenResponse struct {
			AccessToken  string `json:"access_token"`
			RefreshToken string `json:"refresh_token"`
		}

		encodedRefreshToken := usecase.EncodeToken(refreshToken)

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(tokenResponse{AccessToken: accessToken, RefreshToken: encodedRefreshToken})
		if err != nil {
			http.Error(w, "Can't encode json response: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (th *TokenHandler) Refresh() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Provide POST method", http.StatusMethodNotAllowed)
			return
		}
		type tokenRequest struct {
			UserID       string `json:"user_id"`
			RefreshToken string `json:"refresh_token"`
		}
		tokenReq := &tokenRequest{}
		err := json.NewDecoder(r.Body).Decode(tokenReq)
		if err != nil {
			http.Error(w, "Can't decode request, please provide valid data: "+err.Error(), http.StatusBadRequest)
			return
		}

		if tokenReq.RefreshToken == "" || tokenReq.UserID == "" {
			http.Error(w, "Provide valid data", http.StatusBadRequest)
			return
		}

		if err := usecase.ValidateUserID(tokenReq.UserID); err != nil {
			http.Error(w, "Provide valid UUID user_id", http.StatusBadRequest)
			return
		}

		decodedRefreshToken, err := usecase.DecodeToken(tokenReq.RefreshToken)
		if err != nil {
			http.Error(w, "Can't decode refresh token: "+err.Error(), http.StatusInternalServerError)
			return
		}

		hashedToken, err := th.GetTokenByUserID(tokenReq.UserID)
		if err != nil {
			http.Error(w, "Can't get token from db: "+err.Error(), http.StatusInternalServerError)
			return
		}

		err = usecase.CompareTokenAndHash(decodedRefreshToken, hashedToken)
		if err != nil {
			http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
			return
		}

		claims, err := auth.GetRefreshClaims(decodedRefreshToken)
		if err != nil {
			http.Error(w, "Can't parse refresh token: "+err.Error(), http.StatusInternalServerError)
			return
		}

		tokenExpireAt, err := usecase.StrSecToInt(claims.ExpiresAt)
		if err != nil {
			http.Error(w, "Can't parse refresh token: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if ok := time.Now().Before(time.Unix(int64(tokenExpireAt), 0)); !ok {
			http.Error(w, "Refresh token is expired", http.StatusUnauthorized)
			return
		}

		newAccessToken, newRefreshToken, err := auth.GetRefreshedTokens(claims, tokenReq.UserID, r.RemoteAddr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		hashedToken, err = usecase.HashToken(newRefreshToken)
		if err != nil {
			http.Error(w, "Can't hash refresh token: "+err.Error(), http.StatusInternalServerError)
			return
		}

		err = th.RefreshToken(tokenReq.UserID, hashedToken)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		encodedRefreshToken := usecase.EncodeToken(newRefreshToken)

		type tokenResponse struct {
			AccessToken  string `json:"access_token"`
			RefreshToken string `json:"refresh_token"`
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(tokenResponse{AccessToken: newAccessToken, RefreshToken: encodedRefreshToken})
		if err != nil {
			http.Error(w, "Can't encode json response: "+err.Error(), http.StatusInternalServerError)
		}
	}
}
