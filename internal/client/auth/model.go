package auth

type ITokenManager interface {
	LoadTokens() error
	SaveTokens(accessToken, refreshToken string) error
	UpdateTokens(access, refresh string) error
	IsAuthorized() bool
	HandleAuthFailure()
	SetAuthFailCallback(callback func())
}
