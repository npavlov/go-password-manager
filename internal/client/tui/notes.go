package tui

import (
	"context"
	"fmt"

	"github.com/rivo/tview"

	"github.com/npavlov/go-password-manager/internal/client/model"
)

func (t *TUI) showNoteList() {
	list := tview.NewList()

	for _, note := range t.Storage.GetNotes() {
		noteCopy := note
		list.AddItem(note.Content, "(Press Enter to view details)", 0, func() {
			t.showNoteDetails(noteCopy)
		})
	}

	list.AddItem("➕ Add Note", "Create a new note", 'a', func() {
		t.showAddNoteForm()
	})

	list.AddItem("⬅ Back", "Return to main menu", 'b', func() {
		t.App.SetRoot(t.MainMenu(), true)
	})

	list.SetTitle("📝 Notes").SetBorder(true)
	t.App.SetRoot(list, true)
}

func (t *TUI) showNoteDetails(note model.NoteItem) {
	flex := tview.NewFlex().SetDirection(tview.FlexRow)

	textView := tview.NewTextView()
	textView.SetText(fmt.Sprintf(
		"Content: %s\n📌 Metadata:\n",
		note.Content,
	))

	metainfo, err := t.Facade.GetMetainfo(context.Background(), note.ID)
	if err != nil {
		t.Logger.Error().Err(err).Msg("Error getting metainfo")
		textView.SetText(textView.GetText(false) + "\n❌ Error loading metadata")
	} else {
		for key, value := range metainfo {
			textView.SetText(fmt.Sprintf("%s\n%s: %s", textView.GetText(false), key, value))
		}
	}

	textView.SetBorder(true).SetTitle("🔍 Note Details")

	menu := tview.NewList().
		AddItem("🗑 Remove Card", "Delete this note", 'd', func() {
			t.showRemoveNoteForm(note)
		}).
		AddItem("➕ Add Metadata", "Attach metadata to this note", 'm', func() {
			t.showAddMetadataForm(note.StorageItem, func() {
				t.showNoteDetails(note)
			})
		}).
		AddItem("🗑 Remove Metadata", "Delete metadata entry", 'r', func() {
			t.showRemoveMetadataForm(note.StorageItem, func() {
				t.showNoteDetails(note)
			})
		}).
		AddItem("⬅ Back", "Return to note list", 'b', func() {
			t.showNoteList()
		})

	menu.SetBorder(true).SetTitle("⚙ Actions")

	flex.AddItem(textView, 0, 1, false)
	flex.AddItem(menu, 0, 1, true)

	t.App.SetRoot(flex, true)
}

func (t *TUI) showAddNoteForm() {
	form := tview.NewForm()

	form.
		AddTextArea("Content", "", 50, 10, 0, nil).
		AddButton("Save", func() {
			content := form.GetFormItem(0).(*tview.TextArea).GetText()

			noteID, err := t.Facade.StoreNote(context.Background(), content)
			if err != nil {
				t.Logger.Error().Err(err).Msg("Failed to add note")

				return
			}

			// no meta because password is new!
			err = t.Storage.ProcessNote(context.Background(), noteID, map[string]string{})
			if err != nil {
				t.Logger.Error().Err(err).Msg("Failed to add note")

				return
			}

			t.Logger.Info().Msg("Note added successfully")
			t.showNoteList()
		}).
		AddButton("Cancel", func() {
			t.showNoteList()
		}).
		SetTitle("➕ Add New Note").
		SetBorder(true)

	t.App.SetRoot(form, true)
}

func (t *TUI) showRemoveNoteForm(note model.NoteItem) {
	confirmation := tview.NewModal().
		SetText(fmt.Sprintf("Are you sure you want to delete the note titled '%s'?", note.Content)).
		AddButtons([]string{"Yes", "No"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Yes" {
				ok, err := t.Facade.DeleteNote(context.Background(), note.ID)
				if !ok || err != nil {
					t.Logger.Error().Err(err).Msg("Failed to remove note")

					return
				}

				t.Storage.DeleteNotes(note.ID)

				t.Logger.Info().Msg("Note removed successfully")
				t.showNoteList()
			} else {
				t.showNoteDetails(note)
			}
		})

	t.App.SetRoot(confirmation, true)
}
