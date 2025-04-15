package tui_test

import (
	"context"
	"errors"
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"

	"github.com/npavlov/go-password-manager/internal/client/model"
	testutils "github.com/npavlov/go-password-manager/internal/test_utils"
)

func TestShowPasswordList(t *testing.T) {
	ui := setupTUI()

	// Test empty list
	list := ui.ShowPasswordList()
	assert.NotNil(t, list)
	assert.Equal(t, "🔐 Passwords", list.GetTitle())
	assert.Equal(t, 2, list.GetItemCount()) // Only Add and Back buttons

	// Test with items
	mockStorage := ui.Storage.(*testutils.MockStorageManager)
	mockStorage.Passwords = map[string]model.PasswordItem{
		"1": {StorageItem: model.StorageItem{ID: "1"}, Login: "user1"},
		"2": {StorageItem: model.StorageItem{ID: "2"}, Login: "user2"},
	}

	list = ui.ShowPasswordList()
	assert.Equal(t, 4, list.GetItemCount()) // 2 items + Add + Back

	// Test navigation actions
	event := tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone)

	// Test Add Password action
	list.SetCurrentItem(2) // Add Password is third item (after 2 passwords)
	list.InputHandler()(event, nil)

	// Test Back action
	list.SetCurrentItem(3) // Back is last item
	list.InputHandler()(event, nil)
}

func TestShowPasswordDetails(t *testing.T) {
	ui := setupTUI()
	pass := model.PasswordItem{
		StorageItem: model.StorageItem{ID: "123"},
		Login:       "testuser",
		Password:    "secret",
	}

	// Setup mock facade
	mockFacade := ui.Facade.(*testutils.MockFacade)
	mockFacade.GetMetainfoFunc = func(ctx context.Context, id string) (map[string]string, error) {
		return map[string]string{"key": "value"}, nil
	}
	mockFacade.On("GetMetainfo", context.Background(), "123").Return(map[string]string{"key": "value"}, nil)

	flex := ui.ShowPasswordDetails(pass)
	assert.NotNil(t, flex)
	assert.Equal(t, 2, flex.GetItemCount())

	// Verify text view content
	textView := flex.GetItem(0).(*tview.TextView)
	assert.Contains(t, textView.GetText(true), "👤 Username: testuser")
	assert.Contains(t, textView.GetText(true), "🔑 Password: secret")
	assert.Contains(t, textView.GetText(true), "key: value")

	// Verify actions list
	actions := flex.GetItem(1).(*tview.List)
	assert.Equal(t, 5, actions.GetItemCount())
}

func TestShowPasswordDetails_MetaError(t *testing.T) {
	ui := setupTUI()
	pass := model.PasswordItem{
		StorageItem: model.StorageItem{ID: "123"},
		Login:       "testuser",
	}

	// Setup mock facade to return error
	mockFacade := ui.Facade.(*testutils.MockFacade)
	mockFacade.GetMetainfoFunc = func(ctx context.Context, id string) (map[string]string, error) {
		return nil, errors.New("meta error")
	}
	mockFacade.On("GetMetainfo", context.Background(), "123").Return(nil, errors.New("meta error"))

	flex := ui.ShowPasswordDetails(pass)
	textView := flex.GetItem(0).(*tview.TextView)
	assert.Contains(t, textView.GetText(true), "❌ Error loading metadata")
}

func TestShowAddPasswordForm_Success(t *testing.T) {
	ui := setupTUI()

	// Setup mocks
	mockFacade := ui.Facade.(*testutils.MockFacade)
	mockFacade.StorePasswordFunc = func(ctx context.Context, login, password string) (string, error) {
		return "new-id", nil
	}
	mockFacade.On("StorePassword", context.Background(), "testuser", "secret").Return("new-id", nil)

	mockStorage := ui.Storage.(*testutils.MockStorageManager)
	mockStorage.ProcessPasswordFunc = func(ctx context.Context, passID string, meta map[string]string) error {
		assert.Equal(t, "new-id", passID)
		return nil
	}

	form := ui.ShowAddPasswordForm()
	assert.NotNil(t, form)
	assert.Equal(t, "➕ Add New Password", form.GetTitle())

	// Fill form
	form.GetFormItem(0).(*tview.InputField).SetText("testuser")
	form.GetFormItem(1).(*tview.InputField).SetText("secret")

	// Simulate save button click
	event := tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone)
	form.GetButton(0).InputHandler()(event, nil)
}

func TestShowAddPasswordForm_Error(t *testing.T) {
	ui := setupTUI()

	// Setup mock to return error
	mockFacade := ui.Facade.(*testutils.MockFacade)
	mockFacade.StorePasswordFunc = func(ctx context.Context, login, password string) (string, error) {
		return "", errors.New("store failed")
	}
	mockFacade.On("StorePassword", context.Background(), "testuser", "secret").Return("", errors.New("store failed"))

	form := ui.ShowAddPasswordForm()

	// Fill form
	form.GetFormItem(0).(*tview.InputField).SetText("testuser")
	form.GetFormItem(1).(*tview.InputField).SetText("secret")

	// Simulate save button click
	event := tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone)
	form.GetButton(0).InputHandler()(event, nil)
}

func TestShowAddPasswordForm_Cancel(t *testing.T) {
	ui := setupTUI()

	form := ui.ShowAddPasswordForm()

	// Simulate cancel button click
	event := tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone)
	form.GetButton(1).InputHandler()(event, nil)

}

func TestShowChangePasswordForm_Success(t *testing.T) {
	ui := setupTUI()
	pass := model.PasswordItem{
		StorageItem: model.StorageItem{ID: "123"},
		Login:       "testuser",
		Password:    "oldsecret",
	}

	// Setup mocks
	mockFacade := ui.Facade.(*testutils.MockFacade)
	mockFacade.UpdatePasswordFunc = func(ctx context.Context, id, login, password string) error {
		return nil
	}
	mockFacade.On("UpdatePassword", context.Background(), "123", "newuser", "newsecret").Return(nil)

	mockStorage := ui.Storage.(*testutils.MockStorageManager)
	mockStorage.ProcessPasswordFunc = func(ctx context.Context, passID string, meta map[string]string) error {
		assert.Equal(t, "123", passID)
		return nil
	}

	form := ui.ShowChangePasswordForm(pass)
	assert.NotNil(t, form)
	assert.Equal(t, "✏ Change Password", form.GetTitle())

	// Verify pre-populated fields
	assert.Equal(t, "testuser", form.GetFormItem(0).(*tview.InputField).GetText())

	// Update form
	form.GetFormItem(0).(*tview.InputField).SetText("newuser")
	form.GetFormItem(1).(*tview.InputField).SetText("newsecret")

	// Simulate save button click
	event := tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone)
	form.GetButton(0).InputHandler()(event, nil)

}

func TestShowRemovePasswordForm_Confirm(t *testing.T) {
	ui := setupTUI()
	pass := model.PasswordItem{
		StorageItem: model.StorageItem{ID: "123"},
		Login:       "testuser",
	}

	// Setup mocks
	mockFacade := ui.Facade.(*testutils.MockFacade)
	mockFacade.DeletePasswordFunc = func(ctx context.Context, passID string) (bool, error) {
		return true, nil
	}
	mockFacade.On("DeletePassword", context.Background(), "123").Return(true, nil)

	mockStorage := ui.Storage.(*testutils.MockStorageManager)
	mockStorage.DeletePasswordFunc = func(Id string) {
		assert.Equal(t, "123", Id)
	}

	modal := ui.ShowRemovePasswordForm(pass)
	assert.NotNil(t, modal)

	// Simulate "Yes" selection
	modal.SetFocus(0) // Focus "Yes" button
	modal.InputHandler()(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone), nil)

}

func TestShowRemovePasswordForm_Cancel(t *testing.T) {
	ui := setupTUI()
	pass := model.PasswordItem{
		StorageItem: model.StorageItem{ID: "123"},
		Login:       "testuser",
	}

	modal := ui.ShowRemovePasswordForm(pass)

	// Simulate "No" selection
	modal.SetFocus(1) // Focus "No" button
	modal.InputHandler()(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone), nil)

}

func TestShowPasswordDetails_Actions(t *testing.T) {
	ui := setupTUI()
	pass := model.PasswordItem{
		StorageItem: model.StorageItem{
			ID: "123",
			Metadata: map[string]string{
				"key1": "value1",
			},
		},
		Login: "testuser",
	}

	// Setup mock facade
	mockFacade := ui.Facade.(*testutils.MockFacade)
	mockFacade.GetMetainfoFunc = func(ctx context.Context, id string) (map[string]string, error) {
		return pass.Metadata, nil
	}
	mockFacade.On("GetMetainfo", context.Background(), "123").Return(pass.Metadata, nil)

	flex := ui.ShowPasswordDetails(pass)
	actions := flex.GetItem(1).(*tview.List)

	// Test Change Password action
	actions.SetCurrentItem(0) // Change Password is first item
	event := tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone)
	actions.InputHandler()(event, nil)

	// Go back to details
	ui.SetRoot(flex, true)

	// Test Remove Password action
	actions.SetCurrentItem(1) // Remove is second item
	actions.InputHandler()(event, nil)

	// Go back to details
	ui.SetRoot(flex, true)

	// Test Add Metadata action
	actions.SetCurrentItem(2) // Add Metadata is third item
	actions.InputHandler()(event, nil)

	// Go back to details
	ui.SetRoot(flex, true)

	// Test Remove Metadata action
	actions.SetCurrentItem(3) // Remove Metadata is fourth item
	actions.InputHandler()(event, nil)
}
