//nolint:mnd,forcetypeassert
package tui

import (
	"context"
	"fmt"

	"github.com/rivo/tview"

	"github.com/npavlov/go-password-manager/internal/client/model"
)

// ShowPasswordList displays stored passwords.
func (t *TUI) ShowPasswordList() *tview.List {
	list := tview.NewList()

	// List stored passwords
	for _, pass := range t.Storage.GetPasswords() {
		passCopy := pass // Avoid closure issues
		list.AddItem(pass.Login, "(Press Enter to view details)", 0, func() {
			t.SetRoot(t.ShowPasswordDetails(passCopy), true)
		})
	}

	// Option to add a new password
	list.AddItem("➕ Add Password", "Create a new password entry", 'a', func() {
		t.SetRoot(t.ShowAddPasswordForm(), true)
	})

	// Back button
	list.AddItem("⬅ Back", "Return to main menu", 'b', func() { t.App.SetRoot(t.MainMenu(), true) })

	list.SetTitle("🔐 Passwords").SetBorder(true)

	return list
}

// ShowPasswordDetails displays metadata of a selected password.
func (t *TUI) ShowPasswordDetails(pass model.PasswordItem) *tview.Flex {
	flex := tview.NewFlex().SetDirection(tview.FlexRow)

	// Password details view
	textView := tview.NewTextView()
	textView.SetText(fmt.Sprintf("👤 Username: %s\n🔑 Password: %s\n📌 Metadata:\n", pass.Login, pass.Password))

	// Fetch and display metadata
	metainfo, err := t.Facade.GetMetainfo(context.Background(), pass.ID)
	if err != nil {
		t.Logger.Error().Err(err).Msg("Error getting metainfo")
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
			t.SetRoot(t.ShowChangePasswordForm(pass), true)
		}).
		AddItem("🗑Remove Password", "Remove this password", 'c', func() {
			t.SetRoot(t.ShowRemovePasswordForm(pass), true)
		}).
		AddItem("➕ Add Metadata", "Attach new metadata to this password", 'm', func() {
			t.SetRoot(
				t.ShowAddMetadataForm(pass.StorageItem, func() {
					t.SetRoot(t.ShowPasswordDetails(pass), true)
				}),
				true)
		}).
		AddItem("🗑 Remove Metadata", "Delete metadata entry", 'r', func() {
			t.SetRoot(
				t.ShowRemoveMetadataForm(pass.StorageItem, func() {
					t.SetRoot(t.ShowPasswordDetails(pass), true)
				}),
				true,
			)
		}).
		AddItem("⬅ Back", "Return to password list", 'b', func() {
			t.SetRoot(t.ShowPasswordList(), true)
		})

	menu.SetBorder(true).SetTitle("⚙ Actions")

	// Combine views
	flex.AddItem(textView, 0, 1, false)
	flex.AddItem(menu, 0, 1, true)

	return flex
}

// ShowChangePasswordForm allows the user to update a password.
func (t *TUI) ShowChangePasswordForm(pass model.PasswordItem) *tview.Form {
	form := tview.NewForm()

	form.AddInputField("New login", pass.Login, 30, nil, nil).
		AddPasswordField("New Password", "", 30, '*', nil).
		AddButton("Save", func() {
			newLogin := form.GetFormItem(0).(*tview.InputField).GetText()
			newPassword := form.GetFormItem(1).(*tview.InputField).GetText()

			err := t.Facade.UpdatePassword(context.Background(), pass.ID, newLogin, newPassword)
			if err != nil {
				t.Logger.Error().Err(err).Msg("Failed to change password")

				return
			}

			// Update local storage
			err = t.Storage.ProcessPassword(context.Background(), pass.ID, pass.Metadata)
			if err != nil {
				t.Logger.Error().Err(err).Msg("Failed to update password in storage")

				return
			}

			t.Logger.Info().Msg("Password changed successfully")
			t.SetRoot(t.ShowPasswordDetails(pass), true) // Refresh details
		}).
		AddButton("Cancel", func() { t.SetRoot(t.ShowPasswordDetails(pass), true) })

	form.SetTitle("✏ Change Password").SetBorder(true)

	return form
}

// ShowAddPasswordForm displays a form to add a new password entry.
func (t *TUI) ShowAddPasswordForm() *tview.Form {
	form := tview.NewForm()

	form.AddInputField("Username", "", 30, nil, nil).
		AddInputField("Password", "", 30, nil, nil).
		AddButton("Save", func() {
			username := form.GetFormItem(0).(*tview.InputField).GetText()
			password := form.GetFormItem(1).(*tview.InputField).GetText()

			passwordID, err := t.Facade.StorePassword(context.Background(), username, password)
			if err != nil {
				t.Logger.Error().Err(err).Msg("Failed to add password")

				return
			}

			// no meta because password is new!
			err = t.Storage.ProcessPassword(context.Background(), passwordID, map[string]string{})
			if err != nil {
				t.Logger.Error().Err(err).Msg("Failed to add password")

				return
			}

			t.Logger.Info().Msg("Password added successfully")
			t.SetRoot(t.ShowPasswordList(), true) // Refresh list
		}).
		AddButton("Cancel", func() { t.SetRoot(t.ShowPasswordList(), true) })

	form.SetTitle("➕ Add New Password").SetBorder(true)

	return form
}

// ShowRemovePasswordForm displays a confirmation dialog before removing a password.
func (t *TUI) ShowRemovePasswordForm(pass model.PasswordItem) *tview.Modal {
	confirmation := tview.NewModal().
		SetText(fmt.Sprintf("Are you sure you want to delete the password for %s?", pass.Login)).
		AddButtons([]string{"Yes", "No"}).
		SetDoneFunc(func(_ int, buttonLabel string) {
			if buttonLabel == "Yes" {
				ok, err := t.Facade.DeletePassword(context.Background(), pass.ID)
				if !ok || err != nil {
					t.Logger.Error().Err(err).Msg("Failed to remove password")

					return
				}

				t.Storage.DeletePassword(pass.ID)

				t.Logger.Info().Msg("Password removed successfully")
				t.SetRoot(t.ShowPasswordList(), true)
			} else {
				t.SetRoot(t.ShowPasswordDetails(pass), true)
			}
		})

	return confirmation
}
