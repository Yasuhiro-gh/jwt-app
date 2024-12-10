package usecase

type TokenStorage struct {
	store TokenStore
}

func NewTokenStorage(ts TokenStore) *TokenStorage {
	return &TokenStorage{ts}
}

func (t *TokenStorage) SetNewToken(userID, refreshToken string) error {
	return t.store.SetNewToken(userID, refreshToken)
}

func (t *TokenStorage) RefreshToken(userID, refreshToken string) error {
	return t.store.RefreshToken(userID, refreshToken)
}

func (t *TokenStorage) GetTokenByUserID(userID string) (string, error) {
	return t.store.GetTokenByUserID(userID)
}
