package tui

import (
	"context"

	"github.com/npavlov/go-password-manager/internal/client/model"
	"github.com/rivo/tview"
)

// showAddMetadataForm allows the user to add metadata to an item
func (t *TUI) showAddMetadataForm(storageItem model.StorageItem, backFunc func()) {
	form := tview.NewForm()

	form.AddInputField("Key", "", 20, nil, nil).
		AddInputField("Value", "", 50, nil, nil).
		AddButton("Save", func() {
			key := form.GetFormItem(0).(*tview.InputField).GetText()
			value := form.GetFormItem(1).(*tview.InputField).GetText()

			storageItem.Metadata[key] = value

			ok, err := t.facade.SetMetainfo(context.Background(), storageItem.ID, storageItem.Metadata)
			if !ok || err != nil {
				t.logger.Error().Err(err).Msg("Failed to add metadata")

				return
			}

			switch storageItem.Type {
			case model.ItemTypeCard:
				err = t.storage.ProcessCard(context.Background(), storageItem.ID, storageItem.Metadata)
			case model.ItemTypeNote:
				err = t.storage.ProcessNote(context.Background(), storageItem.ID, storageItem.Metadata)
			case model.ItemTypePassword:
				err = t.storage.ProcessPassword(context.Background(), storageItem.ID, storageItem.Metadata)
			}

			if err != nil {
				t.logger.Error().Err(err).Msg("Failed to process metadata")

				return
			}

			t.logger.Info().Msg("Metadata added successfully")
			backFunc() // Refresh details
		}).
		AddButton("Cancel", backFunc)

	form.SetTitle("➕ Add Metadata").SetBorder(true)
	t.app.SetRoot(form, true)
}

// showRemoveMetadataForm allows the user to remove metadata from a password
func (t *TUI) showRemoveMetadataForm(storageItem model.StorageItem, backFunc func()) {
	form := tview.NewForm()

	form.AddInputField("Key", "", 20, nil, nil).
		AddButton("Delete", func() {
			key := form.GetFormItem(0).(*tview.InputField).GetText()

			ok, err := t.facade.DeleteMetainfo(context.Background(), storageItem.ID, key)
			if !ok || err != nil {
				t.logger.Error().Err(err).Msg("Failed to remove metadata")

				return
			}

			delete(storageItem.Metadata, key)

			switch storageItem.Type {
			case model.ItemTypeCard:
				err = t.storage.ProcessCard(context.Background(), storageItem.ID, storageItem.Metadata)
			case model.ItemTypeNote:
				err = t.storage.ProcessNote(context.Background(), storageItem.ID, storageItem.Metadata)
			case model.ItemTypePassword:
				err = t.storage.ProcessPassword(context.Background(), storageItem.ID, storageItem.Metadata)
			}

			if err != nil {
				t.logger.Error().Err(err).Msg("Failed to process metadata")

				return
			}

			t.logger.Info().Msg("Metadata removed successfully")
			backFunc() // Refresh details
		}).
		AddButton("Cancel", backFunc)

	form.SetTitle("🗑 Remove Metadata").SetBorder(true)
	t.app.SetRoot(form, true)
}
