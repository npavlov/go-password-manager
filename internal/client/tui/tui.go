package tui

import (
	"context"

	"github.com/rivo/tview"
	"github.com/rs/zerolog"

	"github.com/npavlov/go-password-manager/internal/client/auth"
	"github.com/npavlov/go-password-manager/internal/client/grpc/facade"
	"github.com/npavlov/go-password-manager/internal/client/storage"
)

type TUI struct {
	App        *tview.Application
	Facade     facade.IFacade
	Storage    storage.IStorageManager
	Logger     *zerolog.Logger
	TokenMgr   auth.ITokenManager
	OnLogin    func()
	OnRegister func()
	stopApp    func()
	SetRoot    func(p tview.Primitive, fullscreen bool) *tview.Application
}

func NewTUI(app *tview.Application, facade facade.IFacade, storage storage.IStorageManager, tokenMgr auth.ITokenManager, log *zerolog.Logger) *TUI {
	return &TUI{
		App:      app,
		Facade:   facade,
		Logger:   log,
		Storage:  storage,
		TokenMgr: tokenMgr,
		SetRoot:  app.SetRoot,
		stopApp:  app.Stop,
	}
}

func (t *TUI) GetApp() *tview.Application {
	return t.SetRoot(t.MainMenu(), true)
}

func (t *TUI) ResetToLoginScreen() {
	t.TokenMgr.HandleAuthFailure()
	t.Logger.Warn().Msg("Authentication failed. Redirecting to login screen...")
	t.SetRoot(t.ShowLoginForm(), true)
}

func (t *TUI) MainMenu() *tview.List {
	menu := tview.NewList()

	if t.TokenMgr.IsAuthorized() {
		menu.AddItem("Passwords", "View and manage stored passwords", 'p', func() {
			t.SetRoot(t.ShowPasswordList(), true)
		})
		menu.AddItem("Notes", "View and manage notes", 'n', func() {
			t.SetRoot(t.ShowNoteList(), true)
		})
		menu.AddItem("Cards", "View and manage cards", 'c', func() {
			t.SetRoot(t.ShowCardList(), true)
		})
		menu.AddItem("Binaries", "View and manage binary files", 'b', func() {
			t.SetRoot(t.ShowBinaryList(), true)
		})
		menu.AddItem("Logout", "Sign out", 'q', t.ResetToLoginScreen)
	} else {
		menu.AddItem("Register", "Create new account", 'r', func() {
			t.SetRoot(t.ShowRegisterForm(), true)
		})
		menu.AddItem("Login", "Sign in", 'l', func() {
			t.SetRoot(t.ShowLoginForm(), true)
		})
		menu.AddItem("Quit", "Exit", 'q', t.stopApp)
	}

	menu.SetTitle("Password Manager").SetBorder(true)

	return menu
}

func (t *TUI) ShowRegisterForm() *tview.Form {
	form := tview.NewForm()

	form.
		AddInputField("Username", "", 20, nil, nil).
		AddPasswordField("Password", "", 20, '*', nil).
		AddInputField("Email", "", 50, nil, nil).
		AddButton("Register", func() { t.HandleRegister(form) }).
		AddButton("Back", func() { t.SetRoot(t.MainMenu(), true) })

	form.SetTitle("Register").SetBorder(true)

	return form
}

func (t *TUI) ShowLoginForm() *tview.Form {
	form := tview.NewForm()

	form.
		AddInputField("Username", "", 20, nil, nil).
		AddPasswordField("Password", "", 20, '*', nil).
		AddButton("Login", func() { t.HandleLogin(form) }).
		AddButton("Back", func() { t.SetRoot(t.MainMenu(), true) })

	form.SetTitle("Login").SetBorder(true)
	return form
}

// ---- Handlers ----

func (t *TUI) HandleRegister(form *tview.Form) {
	username := form.GetFormItem(0).(*tview.InputField).GetText()
	password := form.GetFormItem(1).(*tview.InputField).GetText()
	email := form.GetFormItem(2).(*tview.InputField).GetText()

	userID, err := t.Facade.Register(username, password, email)
	if err != nil {
		t.Logger.Error().Err(err).Msg("Registration failed")
		return
	}

	t.Logger.Info().Str("userID", userID).Msg("Registration successful")

	if t.OnRegister != nil {
		t.OnRegister()
	}

	t.SetRoot(t.MainMenu(), true)
}

func (t *TUI) HandleLogin(form *tview.Form) {
	username := form.GetFormItem(0).(*tview.InputField).GetText()
	password := form.GetFormItem(1).(*tview.InputField).GetText()

	err := t.Facade.Login(username, password)
	if err != nil {
		t.Logger.Error().Err(err).Msg("Login failed")
		return
	}

	if err := t.Storage.SyncItems(context.Background()); err != nil {
		t.Logger.Error().Err(err).Msg("SyncItems failed")
	} else {
		t.Logger.Info().Msg("SyncItems successful")
	}

	if t.OnLogin != nil {
		t.OnLogin()
	}

	t.SetRoot(t.MainMenu(), true)
}
