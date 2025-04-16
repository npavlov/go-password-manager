//nolint:wrapcheck,err113,goconst
package tui_test

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/npavlov/go-password-manager/internal/client/model"
	testutils "github.com/npavlov/go-password-manager/internal/test_utils"
)

func TestShowBinaryList(t *testing.T) {
	t.Parallel()

	ui := setupTUI()

	// Test empty list
	list := ui.ShowBinaryList()
	assert.NotNil(t, list)
	assert.Equal(t, "📦 Binary Files", list.GetTitle())
	assert.Equal(t, 2, list.GetItemCount()) // Only Add and Back buttons

	// Test with items
	mockStorage := ui.Storage.(*testutils.MockStorageManager)
	mockStorage.Binaries = map[string]model.BinaryItem{
		"1": {StorageItem: model.StorageItem{ID: "1"}, Filename: "test1.txt"},
		"2": {StorageItem: model.StorageItem{ID: "2"}, Filename: "test2.txt"},
	}

	list = ui.ShowBinaryList()
	assert.Equal(t, 4, list.GetItemCount()) // 2 items + Add + Back
}

func TestShowBinaryDetails(t *testing.T) {
	t.Parallel()

	ui := setupTUI()
	file := model.BinaryItem{
		StorageItem: model.StorageItem{ID: "123"},
		Filename:    "test.txt",
		Size:        1024,
	}

	details := ui.ShowBinaryDetails(file)
	assert.NotNil(t, details)
	assert.Equal(t, 2, details.GetItemCount())

	// Verify actions list
	actions := details.GetItem(1).(*tview.List)
	assert.Equal(t, 5, actions.GetItemCount())
}

func TestShowBinaryDetails_DownloadAction(t *testing.T) {
	t.Parallel()

	ui := setupTUI()
	file := model.BinaryItem{
		StorageItem: model.StorageItem{ID: "123"},
		Filename:    "test.txt",
		Size:        1024,
	}

	// Setup mock facade
	mockFacade := ui.Facade.(*testutils.MockFacade)
	mockFacade.DownloadBinaryFunc = func(ctx context.Context, fileID string, writer io.Writer) error {
		_, err := writer.Write([]byte("test content"))

		return err
	}
	mockFacade.On("DownloadBinary", t.Context(), "123", mock.Anything).Return(nil)

	details := ui.ShowBinaryDetails(file)
	actions := details.GetItem(1).(*tview.List)

	// Simulate download action
	event := tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone)
	actions.SetCurrentItem(0) // Download is first item
	actions.InputHandler()(event, nil)
}

func TestShowBinaryDetails_RemoveAction(t *testing.T) {
	t.Parallel()

	ui := setupTUI()
	file := model.BinaryItem{
		StorageItem: model.StorageItem{ID: "123"},
		Filename:    "test.txt",
		Size:        1024,
	}

	// Setup mocks
	mockFacade := ui.Facade.(*testutils.MockFacade)
	mockFacade.DeleteBinaryFunc = func(ctx context.Context, fileID string) (bool, error) {
		return true, nil
	}
	mockFacade.On("DeleteBinary", t.Context(), "123").Return(true, nil)

	mockStorage := ui.Storage.(*testutils.MockStorageManager)
	mockStorage.DeleteBinaryFunc = func(Id string) {
		assert.Equal(t, "123", Id)
	}

	details := ui.ShowBinaryDetails(file)
	actions := details.GetItem(1).(*tview.List)

	// Simulate remove action
	event := tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone)
	actions.SetCurrentItem(1) // Remove is second item
	actions.InputHandler()(event, nil)

	// Should show confirmation modal
	modal := ui.ShowRemoveBinaryForm(file)
	ui.SetRoot(modal, true)

	// Simulate confirmation
	modal.InputHandler()(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone), nil)
}

func TestShowUploadBinaryForm_Success(t *testing.T) {
	t.Parallel()

	ui := setupTUI()

	// Create a test file
	tmpFile, err := os.CreateTemp(t.TempDir(), "testfile")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())
	_, _ = tmpFile.WriteString("test content")
	_ = tmpFile.Close()

	// Setup mocks
	mockFacade := ui.Facade.(*testutils.MockFacade)
	mockFacade.UploadBinaryFunc = func(ctx context.Context, filename string, reader io.Reader) (string, error) {
		return "new-id", nil
	}
	mockFacade.On("UploadBinary", t.Context(), filepath.Base(tmpFile.Name()), mock.Anything).Return("new-id", nil)

	mockStorage := ui.Storage.(*testutils.MockStorageManager)
	mockStorage.ProcessBinaryFunc = func(ctx context.Context, fileID string, meta map[string]string) error {
		assert.Equal(t, "new-id", fileID)

		return nil
	}

	form := ui.ShowUploadBinaryForm()
	form.GetFormItem(0).(*tview.InputField).SetText(tmpFile.Name())

	// Simulate upload button click
	event := tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone)
	form.GetButton(0).InputHandler()(event, nil)
}

func TestShowUploadBinaryForm_FileError(t *testing.T) {
	ui := setupTUI()

	form := ui.ShowUploadBinaryForm()
	form.GetFormItem(0).(*tview.InputField).SetText("nonexistent.txt")

	// Simulate upload button click
	event := tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone)
	form.GetButton(0).InputHandler()(event, nil)
}

func TestShowUploadBinaryForm_UploadError(t *testing.T) {
	ui := setupTUI()

	// Create a test file
	tmpFile, err := os.CreateTemp(t.TempDir(), "testfile")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())
	_, _ = tmpFile.WriteString("test content")
	_ = tmpFile.Close()

	// Setup mock to return error
	mockFacade := ui.Facade.(*testutils.MockFacade)
	mockFacade.UploadBinaryFunc = func(ctx context.Context, filename string, reader io.Reader) (string, error) {
		return "", errors.New("upload failed")
	}
	mockFacade.On("UploadBinary", t.Context(), filepath.Base(tmpFile.Name()), mock.Anything).
		Return("", errors.New("upload failed"))

	form := ui.ShowUploadBinaryForm()
	form.GetFormItem(0).(*tview.InputField).SetText(tmpFile.Name())

	// Simulate upload button click
	event := tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone)
	form.GetButton(0).InputHandler()(event, nil)
}

func TestShowRemoveBinaryForm_Confirm(t *testing.T) {
	t.Parallel()

	ui := setupTUI()
	file := model.BinaryItem{
		StorageItem: model.StorageItem{ID: "123"},
		Filename:    "test.txt",
	}

	// Setup mocks
	mockFacade := ui.Facade.(*testutils.MockFacade)
	mockFacade.DeleteBinaryFunc = func(ctx context.Context, fileID string) (bool, error) {
		return true, nil
	}
	mockFacade.On("DeleteBinary", t.Context(), "123").Return(true, nil)

	mockStorage := ui.Storage.(*testutils.MockStorageManager)
	mockStorage.DeleteBinaryFunc = func(Id string) {
		assert.Equal(t, "123", Id)
	}

	modal := ui.ShowRemoveBinaryForm(file)

	// Simulate "Yes" selection
	modal.SetFocus(0) // Focus "Yes" button
	modal.InputHandler()(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone), nil)
}

func TestShowRemoveBinaryForm_Cancel(t *testing.T) {
	ui := setupTUI()
	file := model.BinaryItem{
		StorageItem: model.StorageItem{ID: "123"},
		Filename:    "test.txt",
	}

	modal := ui.ShowRemoveBinaryForm(file)

	// Simulate "No" selection
	modal.SetFocus(1) // Focus "No" button
	modal.InputHandler()(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone), nil)
}

func TestShowBinaryDetails_MetadataActions(t *testing.T) {
	ui := setupTUI()
	file := model.BinaryItem{
		StorageItem: model.StorageItem{
			ID: "123",
			Metadata: map[string]string{
				"key1": "value1",
			},
		},
		Filename: "test.txt",
	}

	details := ui.ShowBinaryDetails(file)
	actions := details.GetItem(1).(*tview.List)

	// Test Add Metadata action
	actions.SetCurrentItem(2) // Add Metadata is third item
	event := tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone)
	actions.InputHandler()(event, nil)

	// Go back to details
	ui.SetRoot(details, true)

	// Test Remove Metadata action
	actions.SetCurrentItem(3) // Remove Metadata is fourth item
	actions.InputHandler()(event, nil)
}
