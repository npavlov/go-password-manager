//nolint:mnd,forcetypeassert,goconst
package tui

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/rivo/tview"

	"github.com/npavlov/go-password-manager/internal/client/model"
)

// ShowBinaryList displays stored binary files.
func (t *TUI) ShowBinaryList() *tview.List {
	list := tview.NewList()

	for _, file := range t.Storage.GetBinaries() {
		binCopy := file
		filename := file.Filename
		list.AddItem(filename, "(Press Enter to view details)", 0, func() {
			t.SetRoot(t.ShowBinaryDetails(binCopy), true)
		})
	}

	list.AddItem("➕ Upload File", "Add a new binary file", 'u', func() {
		t.SetRoot(t.ShowUploadBinaryForm(), true)
	})

	list.AddItem("⬅ Back", "Return to main menu", 'b', func() {
		t.SetRoot(t.MainMenu(), true)
	})

	list.SetTitle("📦 Binary Files").SetBorder(true)

	return list
}

// ShowBinaryDetails displays metadata and allows download.
func (t *TUI) ShowBinaryDetails(file model.BinaryItem) *tview.Flex {
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
				t.Logger.Error().Err(err).Msg("error creating file")
			}
			defer outFile.Close()

			err = t.Facade.DownloadBinary(context.Background(), file.ID, outFile)

			if err != nil {
				t.Logger.Error().Err(err).Msg("Error writing file")
			} else {
				t.Logger.Info().Msgf("File saved to: %s", filePath)

				textView.SetTitle("File saved to: " + file.Filename)
			}
		}).
		AddItem("🗑 Remove Binary", "Delete this file", 'r', func() {
			t.SetRoot(t.ShowRemoveBinaryForm(file), true)
		}).
		AddItem("➕ Add Metadata", "Attach metadata to this note", 'm', func() {
			t.SetRoot(
				t.ShowAddMetadataForm(file.StorageItem, func() {
					t.SetRoot(t.ShowBinaryDetails(file), true)
				}),
				true,
			)
		}).
		AddItem("🗑 Remove Metadata", "Delete metadata entry", 'r', func() {
			t.SetRoot(
				t.ShowRemoveMetadataForm(file.StorageItem, func() {
					t.SetRoot(t.ShowBinaryDetails(file), true)
				}),
				true,
			)
		}).
		AddItem("⬅ Back", "Return to list", 'b', func() {
			t.SetRoot(t.ShowBinaryList(), true)
		})

	menu.SetBorder(true).SetTitle("⚙ Actions")
	flex.AddItem(textView, 0, 1, false)
	flex.AddItem(menu, 0, 1, true)

	return flex
}

// ShowUploadBinaryForm uploads a new file.
func (t *TUI) ShowUploadBinaryForm() *tview.Form {
	form := tview.NewForm()

	form.AddInputField("Local file path", "", 40, nil, nil).
		AddButton("Upload", func() {
			path := form.GetFormItem(0).(*tview.InputField).GetText()

			file, err := os.Open(path)
			if err != nil {
				t.Logger.Error().Err(err).Msg("Failed to open file")

				return
			}
			defer file.Close()

			filename := filepath.Base(path)

			id, err := t.Facade.UploadBinary(context.Background(), filename, file)
			if err != nil {
				t.Logger.Error().Err(err).Msg("Failed to store binary")

				return
			}

			err = t.Storage.ProcessBinary(context.Background(), id, map[string]string{})
			if err != nil {
				t.Logger.Error().Err(err).Msg("Failed to update storage")

				return
			}

			t.Logger.Info().Msgf("Binary uploaded: %s", filename)
			t.SetRoot(t.ShowBinaryList(), true)
		}).
		AddButton("Cancel", func() {
			t.SetRoot(t.ShowBinaryList(), true)
		})

	form.SetTitle("➕ Upload Binary").SetBorder(true)

	return form
}

// ShowRemoveBinaryForm prompts for binary deletion.
func (t *TUI) ShowRemoveBinaryForm(file model.BinaryItem) *tview.Modal {
	confirmation := tview.NewModal().
		SetText(fmt.Sprintf("Delete '%s'?", file.Filename)).
		AddButtons([]string{"Yes", "No"}).
		SetDoneFunc(func(_ int, label string) {
			if label == "Yes" {
				ok, err := t.Facade.DeleteBinary(context.Background(), file.ID)
				if !ok || err != nil {
					t.Logger.Error().Err(err).Msg("Failed to delete binary")

					return
				}

				t.Storage.DeleteBinary(file.ID)

				t.Logger.Info().Msgf("Binary removed: %s", file.Filename)
				t.SetRoot(t.ShowBinaryList(), true)
			} else {
				t.SetRoot(t.ShowBinaryDetails(file), true)
			}
		})

	return confirmation
}
