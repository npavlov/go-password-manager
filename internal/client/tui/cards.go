package tui

import (
	"context"
	"fmt"

	"github.com/rivo/tview"

	"github.com/npavlov/go-password-manager/internal/client/model"
)

// showCardList displays stored cards.
func (t *TUI) showCardList() {
	list := tview.NewList()

	for _, card := range t.storage.GetCards() {
		cardCopy := card

		cardNumber := FormatCardNumber(card.CardNumber)

		list.AddItem(cardNumber, "(Press Enter to view details)", 0, func() {
			t.showCardDetails(cardCopy)
		})
	}

	list.AddItem("➕ Add Card", "Create a new card entry", 'a', func() {
		t.showAddCardForm()
	})

	list.AddItem("⬅ Back", "Return to main menu", 'b', func() {
		t.app.SetRoot(t.mainMenu(), true)
	})

	list.SetTitle("💳 Cards").SetBorder(true)
	t.app.SetRoot(list, true)
}

// showCardDetails displays metadata of a selected card.
func (t *TUI) showCardDetails(card model.CardItem) {
	flex := tview.NewFlex().SetDirection(tview.FlexRow)

	textView := tview.NewTextView()
	textView.SetText(fmt.Sprintf(
		"💳 Card Number: %s\n📅 Expiry: %s\n🔐 CVV: %s\n📌 Metadata:\n",
		card.CardNumber, card.ExpiryDate, card.CVV,
	))

	metainfo, err := t.facade.GetMetainfo(context.Background(), card.ID)
	if err != nil {
		t.logger.Error().Err(err).Msg("Error getting metainfo")
		textView.SetText(textView.GetText(false) + "\n❌ Error loading metadata")
	} else {
		for key, value := range metainfo {
			textView.SetText(fmt.Sprintf("%s\n%s: %s", textView.GetText(false), key, value))
		}
	}

	textView.SetBorder(true).SetTitle("🔍 Card Details")

	menu := tview.NewList().
		AddItem("✏ Edit Card", "Update this card", 'e', func() {
			t.showEditCardForm(card)
		}).
		AddItem("🗑 Remove Card", "Delete this card", 'd', func() {
			t.showRemoveCardForm(card)
		}).
		AddItem("➕ Add Metadata", "Attach metadata to this card", 'm', func() {
			t.showAddMetadataForm(card.StorageItem, func() {
				t.showCardDetails(card)
			})
		}).
		AddItem("🗑 Remove Metadata", "Delete metadata entry", 'r', func() {
			t.showRemoveMetadataForm(card.StorageItem, func() {
				t.showCardDetails(card)
			})
		}).
		AddItem("⬅ Back", "Return to card list", 'b', func() {
			t.showCardList()
		})

	menu.SetBorder(true).SetTitle("⚙ Actions")

	flex.AddItem(textView, 0, 1, false)
	flex.AddItem(menu, 0, 1, true)

	t.app.SetRoot(flex, true)
}

// showAddCardForm creates a form for new card entry.
func (t *TUI) showAddCardForm() {
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

			cardID, err := t.facade.StoreCard(context.Background(), cardNum, expiry, cvv, cardHolder)
			if err != nil {
				t.logger.Error().Err(err).Msg("Failed to add card")

				return
			}

			err = t.storage.ProcessCard(context.Background(), cardID, map[string]string{})
			if err != nil {
				t.logger.Error().Err(err).Msg("Failed to store card")

				return
			}

			t.logger.Info().Msg("Card added successfully")
			t.showCardList()
		}).
		AddButton("Cancel", func() {
			t.showCardList()
		})

	form.SetTitle("➕ Add New Card").SetBorder(true)
	t.app.SetRoot(form, true)
}

// showEditCardForm lets user update card info.
func (t *TUI) showEditCardForm(card model.CardItem) {
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

			err := t.facade.UpdateCard(context.Background(), card.ID, cardNum, expiry, cvv, cardHolder)
			if err != nil {
				t.logger.Error().Err(err).Msg("Failed to update card")

				return
			}

			err = t.storage.ProcessCard(context.Background(), card.ID, card.Metadata)
			if err != nil {
				t.logger.Error().Err(err).Msg("Failed to update local card")

				return
			}

			t.logger.Info().Msg("Card updated successfully")
			t.showCardDetails(card)
		}).
		AddButton("Cancel", func() {
			t.showCardDetails(card)
		})

	form.SetTitle("✏ Edit Card").SetBorder(true)
	t.app.SetRoot(form, true)
}

// showRemoveCardForm confirmation before delete.
func (t *TUI) showRemoveCardForm(card model.CardItem) {
	confirmation := tview.NewModal().
		SetText(fmt.Sprintf("Are you sure you want to delete the card %s?", card.CardNumber)).
		AddButtons([]string{"Yes", "No"}).
		SetDoneFunc(func(buttonIndex int, label string) {
			if label == "Yes" {
				ok, err := t.facade.DeleteCard(context.Background(), card.ID)
				if !ok || err != nil {
					t.logger.Error().Err(err).Msg("Failed to delete card")

					return
				}

				t.storage.DeleteCards(card.ID)
				t.logger.Info().Msg("Card removed successfully")
				t.showCardList()
			} else {
				t.showCardDetails(card)
			}
		})

	t.app.SetRoot(confirmation, true)
}
