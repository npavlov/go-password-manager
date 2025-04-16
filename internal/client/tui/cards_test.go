//nolint:dupl,err113
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

func TestShowCardList(t *testing.T) {
	t.Parallel()

	ui := setupTUI()

	// Test empty list
	list := ui.ShowCardList()
	assert.NotNil(t, list)
	assert.Equal(t, "💳 Cards", list.GetTitle())
	assert.Equal(t, 2, list.GetItemCount()) // Only Add and Back buttons

	// Test with items
	mockStorage := ui.Storage.(*testutils.MockStorageManager)
	mockStorage.Cards = map[string]model.CardItem{
		"1": {StorageItem: model.StorageItem{ID: "1"}, CardNumber: "1111222233334444"},
		"2": {StorageItem: model.StorageItem{ID: "2"}, CardNumber: "5555666677778888"},
	}

	list = ui.ShowCardList()
	assert.Equal(t, 4, list.GetItemCount()) // 2 items + Add + Back

	// Test navigation actions
	event := tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone)

	// Test Add Card action
	list.SetCurrentItem(2) // Add Card is third item (after 2 cards)
	list.InputHandler()(event, nil)

	// Test Back action
	list.SetCurrentItem(3) // Back is last item
	list.InputHandler()(event, nil)
}

func TestShowCardDetails(t *testing.T) {
	t.Parallel()

	ui := setupTUI()
	card := model.CardItem{
		StorageItem: model.StorageItem{ID: "123"},
		CardNumber:  "1111222233334444",
		ExpiryDate:  "12/25",
		CVV:         "123",
	}

	// Setup mock facade
	mockFacade := ui.Facade.(*testutils.MockFacade)
	mockFacade.GetMetainfoFunc = func(ctx context.Context, id string) (map[string]string, error) {
		return map[string]string{"key": "value"}, nil
	}
	mockFacade.On("GetMetainfo", mock.Anything, "123").Return(map[string]string{"key": "value"}, nil)

	textView := ui.ShowCardDetails(card)
	assert.NotNil(t, textView)
	assert.Equal(t, "🔍 Card Details", textView.GetTitle())
	assert.Contains(t, textView.GetText(true), "💳 Card Number: 1111222233334444")
	assert.Contains(t, textView.GetText(true), "📅 Expiry: 12/25")
	assert.Contains(t, textView.GetText(true), "🔐 CVV: 123")
	assert.Contains(t, textView.GetText(true), "key: value")
}

func TestShowCardDetails_MetaError(t *testing.T) {
	t.Parallel()

	ui := setupTUI()
	card := model.CardItem{
		StorageItem: model.StorageItem{ID: "123"},
		CardNumber:  "1111222233334444",
	}

	// Setup mock facade to return error
	mockFacade := ui.Facade.(*testutils.MockFacade)
	mockFacade.GetMetainfoFunc = func(ctx context.Context, id string) (map[string]string, error) {
		return nil, errors.New("meta error")
	}
	mockFacade.On("GetMetainfo", mock.Anything, "123").Return(nil, errors.New("meta error"))

	textView := ui.ShowCardDetails(card)
	assert.Contains(t, textView.GetText(true), "❌ Error loading metadata")
}

func TestShowAddCardForm_Success(t *testing.T) {
	t.Parallel()

	ui := setupTUI()

	// Setup mocks
	mockFacade := ui.Facade.(*testutils.MockFacade)
	mockFacade.StoreCardFunc = func(ctx context.Context, cardNum, expDate, cvv, cardHolder string) (string, error) {
		return "new-id", nil
	}
	mockFacade.On("StoreCard", t.Context(), "1111222233334444", "12/25", "123", "John Doe").Return("new-id", nil)

	mockStorage := ui.Storage.(*testutils.MockStorageManager)
	mockStorage.ProcessCardFunc = func(ctx context.Context, cardID string, meta map[string]string) error {
		assert.Equal(t, "new-id", cardID)

		return nil
	}

	form := ui.ShowAddCardForm()
	assert.NotNil(t, form)
	assert.Equal(t, "➕ Add New Card", form.GetTitle())

	// Fill form
	form.GetFormItem(0).(*tview.InputField).SetText("1111222233334444")
	form.GetFormItem(1).(*tview.InputField).SetText("12/25")
	form.GetFormItem(2).(*tview.InputField).SetText("123")
	form.GetFormItem(3).(*tview.InputField).SetText("John Doe")

	// Simulate save button click
	event := tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone)
	form.GetButton(0).InputHandler()(event, nil)
}

func TestShowAddCardForm_Error(t *testing.T) {
	t.Parallel()

	ui := setupTUI()

	// Setup mock to return error
	mockFacade := ui.Facade.(*testutils.MockFacade)
	mockFacade.StoreCardFunc = func(ctx context.Context, cardNum, expDate, cvv, cardHolder string) (string, error) {
		return "", errors.New("store failed")
	}
	mockFacade.On("StoreCard", t.Context(), "1111222233334444", "12/25", "123", "John Doe").
		Return("", errors.New("store failed"))

	form := ui.ShowAddCardForm()

	// Fill form
	form.GetFormItem(0).(*tview.InputField).SetText("1111222233334444")
	form.GetFormItem(1).(*tview.InputField).SetText("12/25")
	form.GetFormItem(2).(*tview.InputField).SetText("123")
	form.GetFormItem(3).(*tview.InputField).SetText("John Doe")

	// Simulate save button click
	event := tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone)
	form.GetButton(0).InputHandler()(event, nil)
}

func TestShowAddCardForm_Cancel(t *testing.T) {
	ui := setupTUI()

	form := ui.ShowAddCardForm()

	// Simulate cancel button click
	event := tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone)
	form.GetButton(1).InputHandler()(event, nil)
}

func TestShowEditCardForm_Success(t *testing.T) {
	ui := setupTUI()
	card := model.CardItem{
		StorageItem: model.StorageItem{ID: "123"},
		CardNumber:  "1111222233334444",
		ExpiryDate:  "12/25",
		CVV:         "123",
	}

	// Setup mocks
	mockFacade := ui.Facade.(*testutils.MockFacade)
	mockFacade.UpdateCardFunc = func(ctx context.Context, id, cardNum, expDate, cvv, cardHolder string) error {
		return nil
	}
	mockFacade.On("UpdateCard", t.Context(), "123", "4444333322221111", "12/26", "456", "Jane Doe").Return(nil)

	mockStorage := ui.Storage.(*testutils.MockStorageManager)
	mockStorage.ProcessCardFunc = func(ctx context.Context, cardID string, meta map[string]string) error {
		assert.Equal(t, "123", cardID)

		return nil
	}

	form := ui.ShowEditCardForm(card)
	assert.NotNil(t, form)
	assert.Equal(t, "✏ Edit Card", form.GetTitle())

	// Verify pre-populated fields
	assert.Equal(t, "1111222233334444", form.GetFormItem(0).(*tview.InputField).GetText())
	assert.Equal(t, "12/25", form.GetFormItem(1).(*tview.InputField).GetText())
	assert.Equal(t, "123", form.GetFormItem(2).(*tview.InputField).GetText())

	// Update form
	form.GetFormItem(0).(*tview.InputField).SetText("4444333322221111")
	form.GetFormItem(1).(*tview.InputField).SetText("12/26")
	form.GetFormItem(2).(*tview.InputField).SetText("456")
	form.GetFormItem(3).(*tview.InputField).SetText("Jane Doe")

	// Simulate save button click
	event := tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone)
	form.GetButton(0).InputHandler()(event, nil)
}

func TestShowRemoveCardForm_Confirm(t *testing.T) {
	ui := setupTUI()
	card := model.CardItem{
		StorageItem: model.StorageItem{ID: "123"},
		CardNumber:  "1111222233334444",
	}

	// Setup mocks
	mockFacade := ui.Facade.(*testutils.MockFacade)
	mockFacade.DeleteCardFunc = func(ctx context.Context, cardID string) (bool, error) {
		return true, nil
	}
	mockFacade.On("DeleteCard", t.Context(), "123").Return(true, nil)

	mockStorage := ui.Storage.(*testutils.MockStorageManager)
	mockStorage.DeleteCardsFunc = func(Id string) {
		assert.Equal(t, "123", Id)
	}

	modal := ui.ShowRemoveCardForm(card)
	assert.NotNil(t, modal)

	// Simulate "Yes" selection
	modal.SetFocus(0) // Focus "Yes" button
	modal.InputHandler()(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone), nil)
}

func TestShowRemoveCardForm_Cancel(t *testing.T) {
	ui := setupTUI()
	card := model.CardItem{
		StorageItem: model.StorageItem{ID: "123"},
		CardNumber:  "1111222233334444",
	}

	modal := ui.ShowRemoveCardForm(card)

	// Simulate "No" selection
	modal.SetFocus(1) // Focus "No" button
	modal.InputHandler()(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone), nil)
}
