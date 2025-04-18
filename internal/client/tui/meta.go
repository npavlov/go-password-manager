//nolint:mnd,forcetypeassert
package tui

import (
	"context"

	"github.com/rivo/tview"

	"github.com/npavlov/go-password-manager/internal/client/model"
)

// ShowAddMetadataForm allows the user to add metadata to an item.
func (t *TUI) ShowAddMetadataForm(storageItem model.StorageItem, backFunc func()) *tview.Form {
	form := tview.NewForm()

	form.AddInputField("Key", "", 20, nil, nil).
		AddInputField("Value", "", 50, nil, nil).
		AddButton("Save", func() {
			key := form.GetFormItem(0).(*tview.InputField).GetText()
			value := form.GetFormItem(1).(*tview.InputField).GetText()

			storageItem.Metadata[key] = value

			ok, err := t.Facade.SetMetainfo(context.Background(), storageItem.ID, storageItem.Metadata)
			if !ok || err != nil {
				t.Logger.Error().Err(err).Msg("Failed to add metadata")

				return
			}

			switch storageItem.Type {
			case model.ItemTypeCard:
				err = t.Storage.ProcessCard(context.Background(), storageItem.ID, storageItem.Metadata)
			case model.ItemTypeNote:
				err = t.Storage.ProcessNote(context.Background(), storageItem.ID, storageItem.Metadata)
			case model.ItemTypePassword:
				err = t.Storage.ProcessPassword(context.Background(), storageItem.ID, storageItem.Metadata)
			case model.ItemTypeBinary:
				err = t.Storage.ProcessBinary(context.Background(), storageItem.ID, storageItem.Metadata)
			}

			if err != nil {
				t.Logger.Error().Err(err).Msg("Failed to process metadata")

				return
			}

			t.Logger.Info().Msg("Metadata added successfully")
			backFunc() // Refresh details
		}).
		AddButton("Cancel", backFunc)

	form.SetTitle("➕ Add Metadata").SetBorder(true)

	return form
}

// ShowRemoveMetadataForm allows the user to remove metadata from a password.
func (t *TUI) ShowRemoveMetadataForm(storageItem model.StorageItem, backFunc func()) *tview.Form {
	form := tview.NewForm()

	form.AddInputField("Key", "", 20, nil, nil).
		AddButton("Delete", func() {
			key := form.GetFormItem(0).(*tview.InputField).GetText()

			ok, err := t.Facade.DeleteMetainfo(context.Background(), storageItem.ID, key)
			if !ok || err != nil {
				t.Logger.Error().Err(err).Msg("Failed to remove metadata")

				return
			}

			delete(storageItem.Metadata, key)

			switch storageItem.Type {
			case model.ItemTypeCard:
				err = t.Storage.ProcessCard(context.Background(), storageItem.ID, storageItem.Metadata)
			case model.ItemTypeNote:
				err = t.Storage.ProcessNote(context.Background(), storageItem.ID, storageItem.Metadata)
			case model.ItemTypePassword:
				err = t.Storage.ProcessPassword(context.Background(), storageItem.ID, storageItem.Metadata)
			case model.ItemTypeBinary:
				err = t.Storage.ProcessBinary(context.Background(), storageItem.ID, storageItem.Metadata)
			}

			if err != nil {
				t.Logger.Error().Err(err).Msg("Failed to process metadata")

				return
			}

			t.Logger.Info().Msg("Metadata removed successfully")
			backFunc() // Refresh details
		}).
		AddButton("Cancel", backFunc)

	form.SetTitle("🗑 Remove Metadata").SetBorder(true)

	return form
}
