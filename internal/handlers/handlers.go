package handlers

import (
	"encoding/json"
	"github.com/Yasuhiro-gh/jwt-app/internal/auth"
	"github.com/Yasuhiro-gh/jwt-app/internal/usecase"
	"net/http"
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
		userID := r.PathValue("uid")
		if userID == "" {
			http.Error(w, "Provide user id in parameters", http.StatusInternalServerError)
			return
		}
		//Todo: Validate user_id
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

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(tokenResponse{AccessToken: accessToken, RefreshToken: refreshToken})
		if err != nil {
			http.Error(w, "Can't encode json response: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (th *TokenHandler) Refresh() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type tokenResponse struct {
			AccessToken  string `json:"access_token"`
			RefreshToken string `json:"refresh_token"`
		}
		tokens := &tokenResponse{}
		err := json.NewDecoder(r.Body).Decode(tokens)
		if err != nil {
			http.Error(w, "Can't decode refresh token: "+err.Error(), http.StatusInternalServerError)
		}
		claims, err := auth.GetClaimsFromRefreshToken(tokens.RefreshToken)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		hashedToken, err := th.GetTokenByUserID(claims.UserID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = usecase.CompareHashAndToken(tokens.RefreshToken, hashedToken)
		if err != nil {
			http.Error(w, "Invalid refresh token", http.StatusInternalServerError)
			return
		}

		err = auth.ValidateAccessToken(tokens.AccessToken)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		newAccessToken, newRefreshToken, err := auth.GetRefreshedTokens(claims, r.RemoteAddr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		hashedToken, err = usecase.HashToken(newRefreshToken)
		if err != nil {
			http.Error(w, "Can't hash refresh token: "+err.Error(), http.StatusInternalServerError)
			return
		}
		err = th.RefreshToken(claims.UserID, hashedToken)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(tokenResponse{AccessToken: newAccessToken, RefreshToken: newRefreshToken})
		if err != nil {
			http.Error(w, "Can't encode json response: "+err.Error(), http.StatusInternalServerError)
		}
	}
}
