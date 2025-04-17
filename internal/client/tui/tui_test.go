//nolint:err113,forcetypeassert
package tui_test

import (
	"context"
	"errors"
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/npavlov/go-password-manager/internal/client/tui"
	testutils "github.com/npavlov/go-password-manager/internal/test_utils"
)

func setupTUI() *tui.TUI {
	app := tview.NewApplication()
	mockFacade := new(testutils.MockFacade)
	mockStorage := testutils.NewMockStorageManager()
	mockTokenMgr := new(testutils.MockTokenManager)
	logger := testutils.GetTLogger()

	return tui.NewTUI(app, mockFacade, mockStorage, mockTokenMgr, logger)
}

func TestMainMenu_Authorized(t *testing.T) {
	t.Parallel()

	ui := setupTUI()
	mockToken := ui.TokenMgr.(*testutils.MockTokenManager)
	mockToken.Authorized = true
	mockToken.On("IsAuthorized").Return(true)

	menu := ui.MainMenu()

	assert.NotNil(t, menu)
	assert.GreaterOrEqual(t, menu.GetItemCount(), 4)
}

func TestMainMenu_NotAuthorized(t *testing.T) {
	t.Parallel()

	ui := setupTUI()
	mockToken := ui.TokenMgr.(*testutils.MockTokenManager)
	mockToken.Authorized = false
	mockToken.On("IsAuthorized").Return(false)

	menu := ui.MainMenu()

	assert.NotNil(t, menu)
	assert.GreaterOrEqual(t, menu.GetItemCount(), 2)
}

func TestHandleLogin_Success(t *testing.T) {
	t.Parallel()

	ui := setupTUI()

	mockFacade := ui.Facade.(*testutils.MockFacade)
	mockStorage := ui.Storage.(*testutils.MockStorageManager)

	mockStorage.SyncItemsFunc = func(_ context.Context) error {
		return nil
	}

	mockFacade.LoginFunc = func(_, _ string) error {
		return nil
	}

	mockFacade.On("Login", "user", "pass").Return(nil)
	mockStorage.On("SyncItems", mock.Anything).Return(nil)

	form := tview.NewForm().
		AddInputField("Username", "user", 20, nil, nil).
		AddPasswordField("Password", "pass", 20, '*', nil)

	ui.HandleLogin(form)

	mockFacade.AssertExpectations(t)
	mockStorage.AssertExpectations(t)
}

func TestHandleLogin_Failure(t *testing.T) {
	t.Parallel()

	ui := setupTUI()

	mockFacade := ui.Facade.(*testutils.MockFacade)
	mockFacade.LoginFunc = func(_, _ string) error {
		return errors.New("error")
	}
	mockFacade.On("Login", "user", "wrongpass").Return(errors.New("invalid"))

	form := tview.NewForm().
		AddInputField("Username", "user", 20, nil, nil).
		AddPasswordField("Password", "wrongpass", 20, '*', nil)

	ui.HandleLogin(form) // no panic is a pass
}

func TestHandleRegister_Success(t *testing.T) {
	t.Parallel()

	ui := setupTUI()

	mockFacade := ui.Facade.(*testutils.MockFacade)
	mockFacade.RegisterFunc = func(_, _, _ string) (string, error) {
		return "", nil
	}
	mockFacade.On("Register", "user", "pass", "email@example.com").Return("user-id-123", nil)

	form := tview.NewForm().
		AddInputField("Username", "user", 20, nil, nil).
		AddPasswordField("Password", "pass", 20, '*', nil).
		AddInputField("Email", "email@example.com", 50, nil, nil)

	ui.HandleRegister(form)

	mockFacade.AssertExpectations(t)
}

func TestResetToLoginScreen(t *testing.T) {
	t.Parallel()

	ui := setupTUI()
	called := false

	ui.SetRoot = func(_ tview.Primitive, _ bool) *tview.Application {
		called = true

		return ui.App
	}

	mockToken := ui.TokenMgr.(*testutils.MockTokenManager)
	mockToken.On("HandleAuthFailure").Return()

	ui.ResetToLoginScreen()
	assert.True(t, called)
}

func TestShowRegisterForm(t *testing.T) {
	t.Parallel()

	ui := setupTUI()
	mockFacade := ui.Facade.(*testutils.MockFacade)
	mockFacade.RegisterFunc = func(_, _, _ string) (string, error) {
		return "", nil
	}
	mockFacade.On("Register", "testuser", "secret", "test@example.com").Return("user-id-123", nil)

	form := ui.ShowRegisterForm()

	assert.NotNil(t, form)
	assert.Equal(t, "Register", form.GetTitle())

	// Check form field labels
	assert.Equal(t, "Username", form.GetFormItem(0).GetLabel())
	assert.Equal(t, "Password", form.GetFormItem(1).GetLabel())
	assert.Equal(t, "Email", form.GetFormItem(2).GetLabel())

	usernameField := form.GetFormItem(0).(*tview.InputField)
	passwordField := form.GetFormItem(1).(*tview.InputField)
	emailField := form.GetFormItem(2).(*tview.InputField)

	usernameField.SetText("testuser")
	passwordField.SetText("secret")
	emailField.SetText("test@example.com")

	// Trigger the button manually
	event := tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone)
	form.GetButton(0).InputHandler()(event, nil)

	mockFacade.AssertExpectations(t)
}
