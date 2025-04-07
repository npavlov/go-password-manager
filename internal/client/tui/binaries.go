package tui

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/npavlov/go-password-manager/internal/client/model"

	"github.com/rivo/tview"
)

// showBinaryList displays stored binary files
func (t *TUI) showBinaryList() {
	list := tview.NewList()

	for _, file := range t.storage.Binaries {
		binCopy := file
		filename := file.Filename
		list.AddItem(filename, "(Press Enter to view details)", 0, func() {
			t.showBinaryDetails(binCopy)
		})
	}

	list.AddItem("➕ Upload File", "Add a new binary file", 'u', func() {
		t.showUploadBinaryForm()
	})

	list.AddItem("⬅ Back", "Return to main menu", 'b', func() {
		t.app.SetRoot(t.mainMenu(), true)
	})

	list.SetTitle("📦 Binary Files").SetBorder(true)
	t.app.SetRoot(list, true)
}

// showBinaryDetails displays metadata and allows download
func (t *TUI) showBinaryDetails(file model.BinaryItem) {
	flex := tview.NewFlex().SetDirection(tview.FlexRow)

	textView := tview.NewTextView().
		SetText(fmt.Sprintf("📄 Filename: %s\n📦 Size: %d bytes\n", file.Filename, file.Size)).
		SetBorder(true).
		SetTitle("📁 Binary Details")

	// Actions
	menu := tview.NewList().
		AddItem("⬇ Download File", "Save this file to /tmp", 'd', func() {
			filePath := filepath.Join("tmp", file.Filename)
			outFile, err := os.Create(filePath)
			if err != nil {
				t.logger.Error().Err(err).Msg("error creating file")
			}
			defer outFile.Close()

			err = t.facade.DownloadBinary(context.Background(), file.ID, outFile)

			if err != nil {
				t.logger.Error().Err(err).Msg("Error writing file")
			} else {
				t.logger.Info().Msgf("File saved to: %s", filePath)

				textView.SetTitle("File saved to: " + file.Filename)
			}
		}).
		AddItem("🗑 Remove Binary", "Delete this file", 'r', func() {
			t.showRemoveBinaryForm(file)
		}).
		AddItem("➕ Add Metadata", "Attach metadata to this note", 'm', func() {
			t.showAddMetadataForm(file.StorageItem, func() {
				t.showBinaryDetails(file)
			})
		}).
		AddItem("🗑 Remove Metadata", "Delete metadata entry", 'r', func() {
			t.showRemoveMetadataForm(file.StorageItem, func() {
				t.showBinaryDetails(file)
			})
		}).
		AddItem("⬅ Back", "Return to list", 'b', func() {
			t.showBinaryList()
		})

	menu.SetBorder(true).SetTitle("⚙ Actions")
	flex.AddItem(textView, 0, 1, false)
	flex.AddItem(menu, 0, 1, true)

	t.app.SetRoot(flex, true)
}

// showUploadBinaryForm uploads a new file
func (t *TUI) showUploadBinaryForm() {
	form := tview.NewForm()

	form.AddInputField("Local file path", "", 40, nil, nil).
		AddButton("Upload", func() {
			path := form.GetFormItem(0).(*tview.InputField).GetText()

			file, err := os.Open(path)
			if err != nil {
				t.logger.Error().Err(err).Msg("Failed to open file")
				return
			}
			defer file.Close()

			filename := filepath.Base(path)

			id, err := t.facade.UploadBinary(context.Background(), filename, file)
			if err != nil {
				t.logger.Error().Err(err).Msg("Failed to store binary")
				return
			}

			err = t.storage.ProcessBinary(context.Background(), id, map[string]string{})
			if err != nil {
				t.logger.Error().Err(err).Msg("Failed to update storage")
				return
			}

			t.logger.Info().Msgf("Binary uploaded: %s", filename)
			t.showBinaryList()
		}).
		AddButton("Cancel", func() {
			t.showBinaryList()
		})

	form.SetTitle("➕ Upload Binary").SetBorder(true)
	t.app.SetRoot(form, true)
}

// showRemoveBinaryForm prompts for binary deletion
func (t *TUI) showRemoveBinaryForm(file model.BinaryItem) {
	confirmation := tview.NewModal().
		SetText(fmt.Sprintf("Delete '%s'?", file.Filename)).
		AddButtons([]string{"Yes", "No"}).
		SetDoneFunc(func(i int, label string) {
			if label == "Yes" {
				ok, err := t.facade.DeleteBinary(context.Background(), file.ID)
				if !ok || err != nil {
					t.logger.Error().Err(err).Msg("Failed to delete binary")
					return
				}

				delete(t.storage.Binaries, file.ID)
				t.logger.Info().Msgf("Binary removed: %s", file.Filename)
				t.showBinaryList()
			} else {
				t.showBinaryDetails(file)
			}
		})

	t.app.SetRoot(confirmation, true)
}
