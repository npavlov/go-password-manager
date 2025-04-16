//nolint:err113
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

func TestShowAddMetadataForm(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		itemType      model.ItemType
		setupMock     func(*testutils.MockFacade, *testutils.MockStorageManager)
		expectSuccess bool
	}{
		{
			name:     "Success for Card",
			itemType: model.ItemTypeCard,
			setupMock: func(f *testutils.MockFacade, s *testutils.MockStorageManager) {
				f.On("SetMetainfo", t.Context(), "123", mock.Anything).Return(true, nil)
				s.ProcessCardFunc = func(ctx context.Context, id string, meta map[string]string) error {
					return nil
				}
			},
			expectSuccess: true,
		},
		{
			name:     "Success for Note",
			itemType: model.ItemTypeNote,
			setupMock: func(f *testutils.MockFacade, s *testutils.MockStorageManager) {
				f.On("SetMetainfo", t.Context(), "123", mock.Anything).Return(true, nil)
				s.ProcessNoteFunc = func(ctx context.Context, id string, meta map[string]string) error {
					return nil
				}
			},
			expectSuccess: true,
		},
		{
			name:     "Success for Password",
			itemType: model.ItemTypePassword,
			setupMock: func(f *testutils.MockFacade, s *testutils.MockStorageManager) {
				f.On("SetMetainfo", t.Context(), "123", mock.Anything).Return(true, nil)
				s.ProcessPasswordFunc = func(ctx context.Context, id string, meta map[string]string) error {
					return nil
				}
			},
			expectSuccess: true,
		},
		{
			name:     "SetMetainfo Failure",
			itemType: model.ItemTypePassword,
			setupMock: func(f *testutils.MockFacade, s *testutils.MockStorageManager) {
				f.On("SetMetainfo", t.Context(), "123", mock.Anything).Return(false, errors.New("failed"))
			},
			expectSuccess: false,
		},
		{
			name:     "Process Failure",
			itemType: model.ItemTypePassword,
			setupMock: func(f *testutils.MockFacade, s *testutils.MockStorageManager) {
				f.On("SetMetainfo", t.Context(), "123", mock.Anything).Return(true, nil)
				s.ProcessPasswordFunc = func(ctx context.Context, id string, meta map[string]string) error {
					return errors.New("failed")
				}
			},
			expectSuccess: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ui := setupTUI()
			mockFacade := ui.Facade.(*testutils.MockFacade)
			mockStorage := ui.Storage.(*testutils.MockStorageManager)
			tt.setupMock(mockFacade, mockStorage)

			item := model.StorageItem{
				ID:       "123",
				Type:     tt.itemType,
				Metadata: make(map[string]string),
			}

			calledBack := false
			backFunc := func() { calledBack = true }

			form := ui.ShowAddMetadataForm(item, backFunc)
			assert.NotNil(t, form)
			assert.Equal(t, "➕ Add Metadata", form.GetTitle())

			// Fill form
			form.GetFormItem(0).(*tview.InputField).SetText("testkey")
			form.GetFormItem(1).(*tview.InputField).SetText("testvalue")

			if tt.expectSuccess {
				mockFacade.SetMetainfoFunc = func(ctx context.Context, id string, meta map[string]string) (bool, error) {
					return true, nil
				}
			} else {
				mockFacade.SetMetainfoFunc = func(ctx context.Context, id string, meta map[string]string) (bool, error) {
					return false, errors.New("failed")
				}
			}

			// Simulate save button click
			event := tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone)
			form.GetButton(0).InputHandler()(event, nil)

			if tt.expectSuccess {
				assert.True(t, calledBack)
				assert.Equal(t, "testvalue", item.Metadata["testkey"])
			} else {
				assert.False(t, calledBack)
			}

			// Test cancel
			calledBack = false
			form.GetButton(1).InputHandler()(event, nil)
			assert.True(t, calledBack)
		})
	}
}

func TestShowRemoveMetadataForm(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		itemType      model.ItemType
		setupMock     func(*testutils.MockFacade, *testutils.MockStorageManager)
		expectSuccess bool
	}{
		{
			name:     "Success for Card",
			itemType: model.ItemTypeCard,
			setupMock: func(f *testutils.MockFacade, s *testutils.MockStorageManager) {
				f.On("DeleteMetainfo", t.Context(), "123", "testkey").Return(true, nil)
				s.ProcessCardFunc = func(ctx context.Context, id string, meta map[string]string) error {
					return nil
				}
			},
			expectSuccess: true,
		},
		{
			name:     "Success for Note",
			itemType: model.ItemTypeNote,
			setupMock: func(f *testutils.MockFacade, s *testutils.MockStorageManager) {
				f.On("DeleteMetainfo", t.Context(), "123", "testkey").Return(true, nil)
				s.ProcessNoteFunc = func(ctx context.Context, id string, meta map[string]string) error {
					return nil
				}
			},
			expectSuccess: true,
		},
		{
			name:     "Success for Password",
			itemType: model.ItemTypePassword,
			setupMock: func(f *testutils.MockFacade, s *testutils.MockStorageManager) {
				f.On("DeleteMetainfo", t.Context(), "123", "testkey").Return(true, nil)
				s.ProcessPasswordFunc = func(ctx context.Context, id string, meta map[string]string) error {
					return nil
				}
			},
			expectSuccess: true,
		},
		{
			name:     "DeleteMetainfo Failure",
			itemType: model.ItemTypePassword,
			setupMock: func(f *testutils.MockFacade, s *testutils.MockStorageManager) {
				f.On("DeleteMetainfo", t.Context(), "123", "testkey").Return(false, errors.New("failed"))
			},
			expectSuccess: false,
		},
		{
			name:     "Process Failure",
			itemType: model.ItemTypePassword,
			setupMock: func(f *testutils.MockFacade, s *testutils.MockStorageManager) {
				f.On("DeleteMetainfo", t.Context(), "123", "testkey").Return(true, nil)
				s.ProcessPasswordFunc = func(ctx context.Context, id string, meta map[string]string) error {
					return errors.New("failed")
				}
			},
			expectSuccess: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ui := setupTUI()
			mockFacade := ui.Facade.(*testutils.MockFacade)
			mockStorage := ui.Storage.(*testutils.MockStorageManager)
			tt.setupMock(mockFacade, mockStorage)

			item := model.StorageItem{
				ID:   "123",
				Type: tt.itemType,
				Metadata: map[string]string{
					"testkey": "testvalue",
				},
			}

			calledBack := false
			backFunc := func() { calledBack = true }

			form := ui.ShowRemoveMetadataForm(item, backFunc)
			assert.NotNil(t, form)
			assert.Equal(t, "🗑 Remove Metadata", form.GetTitle())

			if tt.expectSuccess {
				mockFacade.DeleteMetainfoFunc = func(ctx context.Context, id, key string) (bool, error) {
					return true, nil
				}
			} else {
				mockFacade.DeleteMetainfoFunc = func(ctx context.Context, id, key string) (bool, error) {
					return false, errors.New("failed")
				}
			}

			// Fill form
			form.GetFormItem(0).(*tview.InputField).SetText("testkey")

			// Simulate delete button click
			event := tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone)
			form.GetButton(0).InputHandler()(event, nil)

			if tt.expectSuccess {
				assert.True(t, calledBack)
				_, exists := item.Metadata["testkey"]
				assert.False(t, exists)
			} else {
				assert.False(t, calledBack)
			}

			// Test cancel
			calledBack = false
			form.GetButton(1).InputHandler()(event, nil)
			assert.True(t, calledBack)
		})
	}
}
