package tui

import (
	"context"
	"log"

	"github.com/rivo/tview"
	"github.com/rs/zerolog"

	"github.com/npavlov/go-password-manager/internal/client/auth"
	"github.com/npavlov/go-password-manager/internal/client/grpc/facade"
	"github.com/npavlov/go-password-manager/internal/client/storage"
)

// TUI handles the terminal user interface.
type TUI struct {
	app      *tview.Application
	facade   facade.IFacade
	storage  storage.IStorageManager
	logger   *zerolog.Logger
	tokenMgr auth.ITokenManager
}

// NewTUI creates a new TUI instance.
func NewTUI(app *tview.Application, facade facade.IFacade, storage storage.IStorageManager, tokenMgr auth.ITokenManager, log *zerolog.Logger) *TUI {
	return &TUI{
		app:      app,
		facade:   facade,
		logger:   log,
		storage:  storage,
		tokenMgr: tokenMgr,
	}
}

// Run starts the TUI.
func (t *TUI) Run() error {
	return t.app.SetRoot(t.mainMenu(), true).Run()
}

// ResetToLoginScreen resets the UI and goes back to the login screen.
func (t *TUI) ResetToLoginScreen() {
	t.logger.Warn().Msg("Authentication failed. Redirecting to login screen...")
	t.showLoginForm()
}

// mainMenu creates the main menu UI.
func (t *TUI) mainMenu() *tview.List {
	menu := tview.NewList()

	if t.tokenMgr.IsAuthorized() {
		menu.AddItem("Passwords", "View and manage stored passwords", 'p', func() { t.showPasswordList() })
		menu.AddItem("Notes", "View and manage notes", 'n', func() { t.showNoteList() })
		menu.AddItem("Cards", "View and manage cards", 'c', func() { t.showCardList() })
		menu.AddItem("Binaries", "View and manage binary files", 'b', func() { t.showBinaryList() })
		menu.AddItem("Logout", "Sign out of the application", 'q', func() { t.ResetToLoginScreen() })
	} else {
		menu.AddItem("Register", "Create a new account", 'r', func() { t.showRegisterForm() })
		menu.AddItem("Login", "Sign in to your account", 'l', func() { t.showLoginForm() })
		menu.AddItem("Quit", "Exit the application", 'q', func() { t.app.Stop() })
	}

	menu.SetTitle("Password Manager").SetBorder(true)

	return menu
}

// showRegisterForm displays the registration form.
func (t *TUI) showRegisterForm() {
	form := tview.NewForm()
	form.AddInputField("Username", "", 20, nil, nil).
		AddPasswordField("Password", "", 20, '*', nil).
		AddInputField("Email", "", 50, nil, nil).
		AddButton("Register", func() {
			username := form.GetFormItem(0).(*tview.InputField).GetText()
			password := form.GetFormItem(1).(*tview.InputField).GetText()
			email := form.GetFormItem(2).(*tview.InputField).GetText()

			userID, err := t.facade.Register(username, password, email)
			if err != nil {
				log.Println("Registration failed:", err)

				return
			}

			log.Println("Registration successful! User ID:", userID)
			t.app.SetRoot(t.mainMenu(), true)
		}).
		AddButton("Back", func() { t.app.SetRoot(t.mainMenu(), true) })

	form.SetTitle("Register").SetBorder(true)
	t.app.SetRoot(form, true)
}

// showLoginForm displays the login form.
func (t *TUI) showLoginForm() {
	form := tview.NewForm()
	form.AddInputField("Username", "", 20, nil, nil).
		AddPasswordField("Password", "", 20, '*', nil).
		AddButton("Login", func() {
			username := form.GetFormItem(0).(*tview.InputField).GetText()
			password := form.GetFormItem(1).(*tview.InputField).GetText()

			err := t.facade.Login(username, password)
			if err != nil {
				log.Println("Login failed:", err)

				return
			}

			err = t.storage.SyncItems(context.Background())
			if err != nil {
				t.logger.Error().Err(err).Msg("SyncItems failed")
			}

			t.logger.Info().Msg("SyncItems successful")
			t.app.SetRoot(t.mainMenu(), true)
		}).
		AddButton("Back", func() { t.app.SetRoot(t.mainMenu(), true) })

	form.SetTitle("Login").SetBorder(true)
	t.app.SetRoot(form, true)
}
