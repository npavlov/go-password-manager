package tui_test

import (
	"context"
	"errors"
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/google/uuid"
	testutils "github.com/npavlov/go-password-manager/internal/test_utils"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"

	"github.com/npavlov/go-password-manager/internal/client/tui"
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
	ui := setupTUI()
	mockToken := ui.TokenMgr.(*testutils.MockTokenManager)
	mockToken.Authorized = true
	mockToken.On("IsAuthorized").Return(true)

	menu := ui.MainMenu()

	assert.NotNil(t, menu)
	assert.GreaterOrEqual(t, menu.GetItemCount(), 4)
}

func TestMainMenu_NotAuthorized(t *testing.T) {
	ui := setupTUI()
	mockToken := ui.TokenMgr.(*testutils.MockTokenManager)
	mockToken.Authorized = false
	mockToken.On("IsAuthorized").Return(false)

	menu := ui.MainMenu()

	assert.NotNil(t, menu)
	assert.GreaterOrEqual(t, menu.GetItemCount(), 2)
}

func TestHandleLogin_Success(t *testing.T) {
	ui := setupTUI()

	mockFacade := ui.Facade.(*testutils.MockFacade)
	mockStorage := ui.Storage.(*testutils.MockStorageManager)

	mockStorage.SyncItemsFunc = func(ctx context.Context) error {
		return nil
	}

	mockFacade.LoginFunc = func(username, password string) error {
		return nil
	}

	mockFacade.On("Login", "user", "pass").Return(nil)
	mockStorage.On("SyncItems", context.Background()).Return(nil)

	form := tview.NewForm().
		AddInputField("Username", "user", 20, nil, nil).
		AddPasswordField("Password", "pass", 20, '*', nil)

	loggedIn := false
	ui.OnLogin = func() { loggedIn = true }

	ui.HandleLogin(form)

	assert.True(t, loggedIn)
}

func TestHandleLogin_Failure(t *testing.T) {
	ui := setupTUI()

	mockFacade := ui.Facade.(*testutils.MockFacade)
	mockFacade.LoginFunc = func(username, password string) error {
		return errors.New("error")
	}
	mockFacade.On("Login", "user", "wrongpass").Return(errors.New("invalid"))

	form := tview.NewForm().
		AddInputField("Username", "user", 20, nil, nil).
		AddPasswordField("Password", "wrongpass", 20, '*', nil)

	ui.HandleLogin(form) // no panic is a pass
}

func TestHandleRegister_Success(t *testing.T) {
	ui := setupTUI()

	mockFacade := ui.Facade.(*testutils.MockFacade)
	mockFacade.RegisterFunc = func(username, password, email string) (string, error) {
		return uuid.New().String(), nil
	}
	mockFacade.On("Register", "user", "pass", "email@example.com").Return("user-id-123", nil)

	form := tview.NewForm().
		AddInputField("Username", "user", 20, nil, nil).
		AddPasswordField("Password", "pass", 20, '*', nil).
		AddInputField("Email", "email@example.com", 50, nil, nil)

	registered := false
	ui.OnRegister = func() { registered = true }

	ui.HandleRegister(form)

	assert.True(t, registered)
}

func TestResetToLoginScreen(t *testing.T) {
	ui := setupTUI()
	called := false

	ui.SetRoot = func(p tview.Primitive, fullscreen bool) *tview.Application {
		called = true
		return ui.App
	}

	mockToken := ui.TokenMgr.(*testutils.MockTokenManager)
	mockToken.On("HandleAuthFailure").Return()

	ui.ResetToLoginScreen()
	assert.True(t, called)
}

func TestShowRegisterForm(t *testing.T) {
	ui := setupTUI()
	mockFacade := ui.Facade.(*testutils.MockFacade)
	mockFacade.RegisterFunc = func(username, password, email string) (string, error) {
		return uuid.New().String(), nil
	}

	form := ui.ShowRegisterForm()

	assert.NotNil(t, form)
	assert.Equal(t, "Register", form.GetTitle())

	// Check form field labels
	assert.Equal(t, "Username", form.GetFormItem(0).GetLabel())
	assert.Equal(t, "Password", form.GetFormItem(1).GetLabel())
	assert.Equal(t, "Email", form.GetFormItem(2).GetLabel())

	// Simulate clicking the "Register" button (item 3)
	called := false
	ui.OnRegister = func() { called = true }

	usernameField := form.GetFormItem(0).(*tview.InputField)
	passwordField := form.GetFormItem(1).(*tview.InputField)
	emailField := form.GetFormItem(2).(*tview.InputField)

	usernameField.SetText("testuser")
	passwordField.SetText("secret")
	emailField.SetText("test@example.com")

	// Trigger the button manually
	event := tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone)
	form.GetButton(0).InputHandler()(event, nil)

	assert.True(t, called)
}
