package usecase

type TokenStore interface {
	SetNewToken(string, string) error
	RefreshToken(string) error
	GetTokenByUserID(string) (string, error)
}
