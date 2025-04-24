//nolint:dupl,err113,exhaustruct,forcetypeassert
package tui_test

import (
	"context"
	"errors"
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/npavlov/go-password-manager/internal/client/model"
	testutils "github.com/npavlov/go-password-manager/internal/test_utils"
)

const (
	newID = "new-id"
)

func TestShowNoteList(t *testing.T) {
	t.Parallel()

	ui := setupTUI()

	// Test empty list
	list := ui.ShowNoteList()
	assert.NotNil(t, list)
	assert.Equal(t, "📝 Notes", list.GetTitle())
	assert.Equal(t, 2, list.GetItemCount()) // Only Add and Back buttons

	// Test with items
	mockStorage := ui.Storage.(*testutils.MockStorageManager)
	mockStorage.Notes = map[string]model.NoteItem{
		"1": {StorageItem: model.StorageItem{ID: "1"}, Content: "Note 1"},
		"2": {StorageItem: model.StorageItem{ID: "2"}, Content: "Note 2"},
	}

	list = ui.ShowNoteList()
	assert.Equal(t, 4, list.GetItemCount()) // 2 items + Add + Back

	// Test navigation actions
	event := tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone)

	// Test Add Note action
	list.SetCurrentItem(2) // Add Note is third item (after 2 notes)
	list.InputHandler()(event, nil)

	// Test Back action
	list.SetCurrentItem(3) // Back is last item
	list.InputHandler()(event, nil)
}

func TestShowNoteDetails(t *testing.T) {
	t.Parallel()

	ui := setupTUI()
	note := model.NoteItem{
		StorageItem: model.StorageItem{ID: "123"},
		Content:     "Test note content",
	}

	// Setup mock facade
	mockFacade := ui.Facade.(*testutils.MockFacade)
	mockFacade.GetMetainfoFunc = func(_ context.Context, _ string) (map[string]string, error) {
		return map[string]string{"key": "value"}, nil
	}
	mockFacade.On("GetMetainfo", mock.Anything, "123").Return(map[string]string{"key": "value"}, nil)

	flex := ui.ShowNoteDetails(note)
	assert.NotNil(t, flex)
	assert.Equal(t, 2, flex.GetItemCount())

	// Verify text view content
	textView := flex.GetItem(0).(*tview.TextView)
	assert.Contains(t, textView.GetText(true), "Content: Test note content")
	assert.Contains(t, textView.GetText(true), "key: value")

	// Verify actions list
	actions := flex.GetItem(1).(*tview.List)
	assert.Equal(t, 4, actions.GetItemCount())
}

func TestShowNoteDetails_MetaError(t *testing.T) {
	t.Parallel()

	ui := setupTUI()
	note := model.NoteItem{
		StorageItem: model.StorageItem{ID: "123"},
		Content:     "Test note content",
	}

	// Setup mock facade to return error
	mockFacade := ui.Facade.(*testutils.MockFacade)
	mockFacade.GetMetainfoFunc = func(_ context.Context, _ string) (map[string]string, error) {
		return nil, errors.New("meta error")
	}
	mockFacade.On("GetMetainfo", mock.Anything, "123").Return(nil, errors.New("meta error"))

	flex := ui.ShowNoteDetails(note)
	textView := flex.GetItem(0).(*tview.TextView)
	assert.Contains(t, textView.GetText(true), "❌ Error loading metadata")
}

func TestShowAddNoteForm_Success(t *testing.T) {
	t.Parallel()

	ui := setupTUI()

	// Setup mocks
	mockFacade := ui.Facade.(*testutils.MockFacade)
	mockFacade.StoreNoteFunc = func(_ context.Context, _ string) (string, error) {
		return newID, nil
	}
	mockFacade.On("StoreNote", t.Context(), "test content").Return(newID, nil)

	mockStorage := ui.Storage.(*testutils.MockStorageManager)
	mockStorage.ProcessNoteFunc = func(_ context.Context, noteID string, _ map[string]string) error {
		assert.Equal(t, newID, noteID)

		return nil
	}

	form := ui.ShowAddNoteForm()
	assert.NotNil(t, form)
	assert.Equal(t, "➕ Add New Note", form.GetTitle())

	// Fill form
	textArea := form.GetFormItem(0).(*tview.TextArea)
	textArea.SetText("test content", true)

	// Simulate save button click
	event := tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone)
	form.GetButton(0).InputHandler()(event, nil)
}

func TestShowAddNoteForm_Error(t *testing.T) {
	t.Parallel()

	ui := setupTUI()

	// Setup mock to return error
	mockFacade := ui.Facade.(*testutils.MockFacade)
	mockFacade.StoreNoteFunc = func(_ context.Context, _ string) (string, error) {
		return "", errors.New("store failed")
	}
	mockFacade.On("StoreNote", t.Context(), "test content").Return("", errors.New("store failed"))

	form := ui.ShowAddNoteForm()

	// Fill form
	textArea := form.GetFormItem(0).(*tview.TextArea)
	textArea.SetText("test content", true)

	// Simulate save button click
	event := tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone)
	form.GetButton(0).InputHandler()(event, nil)
}

func TestShowAddNoteForm_Cancel(t *testing.T) {
	t.Parallel()

	ui := setupTUI()

	form := ui.ShowAddNoteForm()

	// Simulate cancel button click
	event := tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone)
	form.GetButton(1).InputHandler()(event, nil)
}

func TestShowRemoveNoteForm_Confirm(t *testing.T) {
	t.Parallel()

	ui := setupTUI()
	note := model.NoteItem{
		StorageItem: model.StorageItem{ID: "123"},
		Content:     "Test note",
	}

	// Setup mocks
	mockFacade := ui.Facade.(*testutils.MockFacade)
	mockFacade.DeleteNoteFunc = func(_ context.Context, _ string) (bool, error) {
		return true, nil
	}
	mockFacade.On("DeleteNote", t.Context(), "123").Return(true, nil)

	mockStorage := ui.Storage.(*testutils.MockStorageManager)
	mockStorage.DeleteNotesFunc = func(Id string) {
		assert.Equal(t, "123", Id)
	}

	modal := ui.ShowRemoveNoteForm(note)
	assert.NotNil(t, modal)

	// Simulate "Yes" selection
	modal.SetFocus(0) // Focus "Yes" button
	modal.InputHandler()(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone), nil)
}

func TestShowRemoveNoteForm_Cancel(t *testing.T) {
	t.Parallel()

	ui := setupTUI()
	note := model.NoteItem{
		StorageItem: model.StorageItem{ID: "123"},
		Content:     "Test note",
	}

	modal := ui.ShowRemoveNoteForm(note)

	// Simulate "No" selection
	modal.SetFocus(1) // Focus "No" button
	modal.InputHandler()(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone), nil)
}

func TestShowNoteDetails_Actions(t *testing.T) {
	t.Parallel()

	ui := setupTUI()
	note := model.NoteItem{
		StorageItem: model.StorageItem{
			ID: "123",
			Metadata: map[string]string{
				"key1": "value1",
			},
		},
		Content: "Test note",
	}

	// Setup mock facade
	mockFacade := ui.Facade.(*testutils.MockFacade)
	mockFacade.On("GetMetainfo", mock.Anything, "123").Return(note.Metadata, nil)
	mockFacade.GetMetainfoFunc = func(_ context.Context, _ string) (map[string]string, error) {
		return note.Metadata, nil
	}

	flex := ui.ShowNoteDetails(note)
	actions := flex.GetItem(1).(*tview.List)

	// Test Remove action
	actions.SetCurrentItem(0) // Remove is first item
	event := tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone)
	actions.InputHandler()(event, nil)

	// Go back to details
	ui.SetRoot(flex, true)

	// Test Add Metadata action
	actions.SetCurrentItem(1) // Add Metadata is second item
	actions.InputHandler()(event, nil)

	// Go back to details
	ui.SetRoot(flex, true)

	// Test Remove Metadata action
	actions.SetCurrentItem(2) // Remove Metadata is third item
	actions.InputHandler()(event, nil)
}
