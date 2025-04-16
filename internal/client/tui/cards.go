//nolint:mnd,forcetypeassert
package tui

import (
	"context"
	"fmt"

	"github.com/rivo/tview"

	"github.com/npavlov/go-password-manager/internal/client/model"
)

// ShowCardList displays stored cards.
func (t *TUI) ShowCardList() *tview.List {
	list := tview.NewList()

	for _, card := range t.Storage.GetCards() {
		cardCopy := card

		cardNumber := FormatCardNumber(card.CardNumber)

		list.AddItem(cardNumber, "(Press Enter to view details)", 0, func() {
			t.App.SetRoot(t.ShowCardDetails(cardCopy), true)
		})
	}

	list.AddItem("➕ Add Card", "Create a new card entry", 'a', func() {
		t.App.SetRoot(t.ShowAddCardForm(), true)
	})

	list.AddItem("⬅ Back", "Return to main menu", 'b', func() {
		t.App.SetRoot(t.MainMenu(), true)
	})

	list.SetTitle("💳 Cards").SetBorder(true)

	return list
}

// ShowCardDetails displays metadata of a selected card.
func (t *TUI) ShowCardDetails(card model.CardItem) *tview.TextView {
	flex := tview.NewFlex().SetDirection(tview.FlexRow)

	textView := tview.NewTextView()
	textView.SetText(fmt.Sprintf(
		"💳 Card Number: %s\n📅 Expiry: %s\n🔐 CVV: %s\n📌 Metadata:\n",
		card.CardNumber, card.ExpiryDate, card.CVV,
	))

	metainfo, err := t.Facade.GetMetainfo(context.Background(), card.ID)
	if err != nil {
		t.Logger.Error().Err(err).Msg("Error getting metainfo")
		textView.SetText(textView.GetText(false) + "\n❌ Error loading metadata")
	} else {
		for key, value := range metainfo {
			textView.SetText(fmt.Sprintf("%s\n%s: %s", textView.GetText(false), key, value))
		}
	}

	textView.SetBorder(true).SetTitle("🔍 Card Details")

	menu := tview.NewList().
		AddItem("✏ Edit Card", "Update this card", 'e', func() {
			t.SetRoot(t.ShowEditCardForm(card), true)
		}).
		AddItem("🗑 Remove Card", "Delete this card", 'd', func() {
			t.SetRoot(t.ShowRemoveCardForm(card), true)
		}).
		AddItem("➕ Add Metadata", "Attach metadata to this card", 'm', func() {
			t.SetRoot(
				t.ShowAddMetadataForm(card.StorageItem, func() {
					t.SetRoot(t.ShowCardDetails(card), true)
				}),
				true,
			)
		}).
		AddItem("🗑 Remove Metadata", "Delete metadata entry", 'r', func() {
			t.SetRoot(
				t.ShowRemoveMetadataForm(card.StorageItem, func() {
					t.SetRoot(t.ShowCardDetails(card), true)
				}), true,
			)
		}).
		AddItem("⬅ Back", "Return to card list", 'b', func() {
			t.SetRoot(t.ShowCardList(), true)
		})

	menu.SetBorder(true).SetTitle("⚙ Actions")

	flex.AddItem(textView, 0, 1, false)
	flex.AddItem(menu, 0, 1, true)

	return textView
}

// ShowAddCardForm creates a form for new card entry.
func (t *TUI) ShowAddCardForm() *tview.Form {
	form := tview.NewForm()

	form.AddInputField("Card Number", "", 30, nil, nil).
		AddInputField("Expiry Date", "", 10, nil, nil).
		AddInputField("CVV", "", 10, nil, nil).
		AddInputField("Card Holder", "", 10, nil, nil).
		AddButton("Save", func() {
			cardNum := form.GetFormItem(0).(*tview.InputField).GetText()
			expiry := form.GetFormItem(1).(*tview.InputField).GetText()
			cvv := form.GetFormItem(2).(*tview.InputField).GetText()
			cardHolder := form.GetFormItem(3).(*tview.InputField).GetText()

			cardID, err := t.Facade.StoreCard(context.Background(), cardNum, expiry, cvv, cardHolder)
			if err != nil {
				t.Logger.Error().Err(err).Msg("Failed to add card")

				return
			}

			err = t.Storage.ProcessCard(context.Background(), cardID, map[string]string{})
			if err != nil {
				t.Logger.Error().Err(err).Msg("Failed to store card")

				return
			}

			t.Logger.Info().Msg("Card added successfully")
			t.SetRoot(t.ShowCardList(), true)
		}).
		AddButton("Cancel", func() {
			t.SetRoot(t.ShowCardList(), true)
		})

	form.SetTitle("➕ Add New Card").SetBorder(true)

	return form
}

// ShowEditCardForm lets user update card info.
func (t *TUI) ShowEditCardForm(card model.CardItem) *tview.Form {
	form := tview.NewForm()

	form.AddInputField("Card Number", card.CardNumber, 30, nil, nil).
		AddInputField("Expiry Date", card.ExpiryDate, 10, nil, nil).
		AddInputField("CVV", card.CVV, 10, nil, nil).
		AddInputField("Card Holder", "", 10, nil, nil).
		AddButton("Save", func() {
			cardNum := form.GetFormItem(0).(*tview.InputField).GetText()
			expiry := form.GetFormItem(1).(*tview.InputField).GetText()
			cvv := form.GetFormItem(2).(*tview.InputField).GetText()
			cardHolder := form.GetFormItem(3).(*tview.InputField).GetText()

			err := t.Facade.UpdateCard(context.Background(), card.ID, cardNum, expiry, cvv, cardHolder)
			if err != nil {
				t.Logger.Error().Err(err).Msg("Failed to update card")

				return
			}

			err = t.Storage.ProcessCard(context.Background(), card.ID, card.Metadata)
			if err != nil {
				t.Logger.Error().Err(err).Msg("Failed to update local card")

				return
			}

			t.Logger.Info().Msg("Card updated successfully")
			t.App.SetRoot(t.ShowCardDetails(card), true)
		}).
		AddButton("Cancel", func() {
			t.SetRoot(t.ShowCardDetails(card), true)
		})

	form.SetTitle("✏ Edit Card").SetBorder(true)

	return form
}

// ShowRemoveCardForm confirmation before delete.
func (t *TUI) ShowRemoveCardForm(card model.CardItem) *tview.Modal {
	confirmation := tview.NewModal().
		SetText(fmt.Sprintf("Are you sure you want to delete the card %s?", card.CardNumber)).
		AddButtons([]string{"Yes", "No"}).
		SetDoneFunc(func(buttonIndex int, label string) {
			if label == "Yes" {
				ok, err := t.Facade.DeleteCard(context.Background(), card.ID)
				if !ok || err != nil {
					t.Logger.Error().Err(err).Msg("Failed to delete card")

					return
				}

				t.Storage.DeleteCards(card.ID)
				t.Logger.Info().Msg("Card removed successfully")
				t.SetRoot(t.ShowCardList(), true)
			} else {
				t.SetRoot(t.ShowCardDetails(card), true)
			}
		})

	return confirmation
}
