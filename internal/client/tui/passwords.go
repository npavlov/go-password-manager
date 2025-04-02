package tui

import (
	"context"
	"fmt"

	"github.com/npavlov/go-password-manager/internal/client/model"
	"github.com/rivo/tview"
)

// showPasswordList displays stored passwords
func (t *TUI) showPasswordList() {
	list := tview.NewList()

	// List stored passwords
	for _, pass := range t.storage.Password {
		passCopy := pass // Avoid closure issues
		list.AddItem(pass.Login, "(Press Enter to view details)", 0, func() {
			t.showPasswordDetails(passCopy)
		})
	}

	// Option to add a new password
	list.AddItem("➕ Add Password", "Create a new password entry", 'a', func() {
		t.showAddPasswordForm()
	})

	// Back button
	list.AddItem("⬅ Back", "Return to main menu", 'b', func() { t.app.SetRoot(t.mainMenu(), true) })

	list.SetTitle("🔐 Passwords").SetBorder(true)
	t.app.SetRoot(list, true)
}

// showPasswordDetails displays metadata of a selected password
func (t *TUI) showPasswordDetails(pass model.PasswordItem) {
	flex := tview.NewFlex().SetDirection(tview.FlexRow)

	// Password details view
	textView := tview.NewTextView()
	textView.SetText(fmt.Sprintf("👤 Username: %s\n🔑 Password: %s\n📌 Metadata:\n", pass.Login, pass.Password))

	// Fetch and display metadata
	metainfo, err := t.facade.GetMetainfo(context.Background(), pass.ID)
	if err != nil {
		t.logger.Error().Err(err).Msg("Error getting metainfo")
		textView.SetText(textView.GetText(false) + "\n❌ Error loading metadata")
	} else {
		for key, value := range metainfo {
			textView.SetText(fmt.Sprintf("%s\n%s: %s", textView.GetText(false), key, value))
		}
	}

	textView.SetBorder(true).SetTitle("🔍 Password Details")

	// Action menu for metadata
	menu := tview.NewList().
		AddItem("✏ Change Password", "Update this password", 'c', func() {
			t.showChangePasswordForm(pass)
		}).
		AddItem("➕ Add Metadata", "Attach new metadata to this password", 'm', func() {
			t.showAddMetadataForm(pass)
		}).
		AddItem("🗑 Remove Metadata", "Delete metadata entry", 'r', func() {
			t.showRemoveMetadataForm(pass)
		}).
		AddItem("⬅ Back", "Return to password list", 'b', func() { t.showPasswordList() })

	menu.SetBorder(true).SetTitle("⚙ Actions")

	// Combine views
	flex.AddItem(textView, 0, 1, false)
	flex.AddItem(menu, 0, 1, true)

	t.app.SetRoot(flex, true)
}

// showChangePasswordForm allows the user to update a password
func (t *TUI) showChangePasswordForm(pass model.PasswordItem) {
	form := tview.NewForm()

	form.AddInputField("New login", pass.Login, 30, nil, nil).
		AddPasswordField("New Password", "", 30, '*', nil).
		AddButton("Save", func() {
			newLogin := form.GetFormItem(0).(*tview.InputField).GetText()
			newPassword := form.GetFormItem(1).(*tview.InputField).GetText()

			err := t.facade.UpdatePassword(context.Background(), pass.ID, newLogin, newPassword)
			if err != nil {
				t.logger.Error().Err(err).Msg("Failed to change password")

				return
			}

			// Update local storage
			err = t.storage.ProcessPassword(context.Background(), pass.ID, pass.Metadata)
			if err != nil {
				t.logger.Error().Err(err).Msg("Failed to update password in storage")

				return
			}

			t.logger.Info().Msg("Password changed successfully")
			t.showPasswordDetails(pass) // Refresh details
		}).
		AddButton("Cancel", func() { t.showPasswordDetails(pass) })

	form.SetTitle("✏ Change Password").SetBorder(true)
	t.app.SetRoot(form, true)
}

// showAddMetadataForm allows the user to add metadata to a password
func (t *TUI) showAddMetadataForm(pass model.PasswordItem) {
	form := tview.NewForm()

	form.AddInputField("Key", "", 20, nil, nil).
		AddInputField("Value", "", 50, nil, nil).
		AddButton("Save", func() {
			key := form.GetFormItem(0).(*tview.InputField).GetText()
			value := form.GetFormItem(1).(*tview.InputField).GetText()

			pass.Metadata[key] = value

			ok, err := t.facade.SetMetainfo(context.Background(), pass.ID, pass.Metadata)
			if !ok || err != nil {
				t.logger.Error().Err(err).Msg("Failed to add metadata")
				return
			}

			err = t.storage.ProcessPassword(context.Background(), pass.ID, pass.Metadata)
			if err != nil {
				t.logger.Error().Err(err).Msg("Failed to add password")

				return
			}

			t.logger.Info().Msg("Metadata added successfully")
			t.showPasswordDetails(pass) // Refresh details
		}).
		AddButton("Cancel", func() { t.showPasswordDetails(pass) })

	form.SetTitle("➕ Add Metadata").SetBorder(true)
	t.app.SetRoot(form, true)
}

// showRemoveMetadataForm allows the user to remove metadata from a password
func (t *TUI) showRemoveMetadataForm(pass model.PasswordItem) {
	form := tview.NewForm()

	form.AddInputField("Key", "", 20, nil, nil).
		AddButton("Delete", func() {
			key := form.GetFormItem(0).(*tview.InputField).GetText()

			ok, err := t.facade.DeleteMetainfo(context.Background(), pass.ID, key)
			if !ok || err != nil {
				t.logger.Error().Err(err).Msg("Failed to remove metadata")

				return
			}

			delete(pass.Metadata, key)

			err = t.storage.ProcessPassword(context.Background(), pass.ID, pass.Metadata)
			if err != nil {
				t.logger.Error().Err(err).Msg("Failed to add password")

				return
			}

			t.logger.Info().Msg("Metadata removed successfully")
			t.showPasswordDetails(pass) // Refresh details
		}).
		AddButton("Cancel", func() { t.showPasswordDetails(pass) })

	form.SetTitle("🗑 Remove Metadata").SetBorder(true)
	t.app.SetRoot(form, true)
}

// showAddPasswordForm displays a form to add a new password entry
func (t *TUI) showAddPasswordForm() {
	form := tview.NewForm()

	form.AddInputField("Username", "", 30, nil, nil).
		AddInputField("Password", "", 30, nil, nil).
		AddButton("Save", func() {
			username := form.GetFormItem(0).(*tview.InputField).GetText()
			password := form.GetFormItem(1).(*tview.InputField).GetText()

			passwordId, err := t.facade.StorePassword(context.Background(), username, password)
			if err != nil {
				t.logger.Error().Err(err).Msg("Failed to add password")
				return
			}

			// no meta because password is new!
			err = t.storage.ProcessPassword(context.Background(), passwordId, map[string]string{})
			if err != nil {
				t.logger.Error().Err(err).Msg("Failed to add password")

				return
			}

			t.logger.Info().Msg("Password added successfully")
			t.showPasswordList() // Refresh list
		}).
		AddButton("Cancel", func() { t.showPasswordList() })

	form.SetTitle("➕ Add New Password").SetBorder(true)
	t.app.SetRoot(form, true)
}
