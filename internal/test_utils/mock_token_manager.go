//nolint:wrapcheck
package testutils

import "github.com/stretchr/testify/mock"

type MockTokenManager struct {
	mock.Mock
	Authorized bool
}

func (m *MockTokenManager) GetAccessToken() string {
	args := m.Called()

	return args.String(0)
}

func (m *MockTokenManager) GetRefreshToken() string {
	args := m.Called()

	return args.String(0)
}

func (m *MockTokenManager) LoadTokens() error {
	return m.Called().Error(0)
}

func (m *MockTokenManager) SaveTokens(accessToken, refreshToken string) error {
	args := m.Called(accessToken, refreshToken)

	return args.Error(0)
}

func (m *MockTokenManager) UpdateTokens(access, refresh string) error {
	args := m.Called(access, refresh)

	return args.Error(0)
}

func (m *MockTokenManager) IsAuthorized() bool {
	return m.Authorized
}

func (m *MockTokenManager) HandleAuthFailure() {
	m.Authorized = false
	m.Called()
}

func (m *MockTokenManager) SetAuthFailCallback(callback func()) {
	m.Called(callback)
}
