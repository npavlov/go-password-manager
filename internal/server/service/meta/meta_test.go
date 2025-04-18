//nolint:exhaustruct
package meta_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	pb "github.com/npavlov/go-password-manager/gen/proto/metadata"
	"github.com/npavlov/go-password-manager/internal/server/config"
	"github.com/npavlov/go-password-manager/internal/server/db"
	"github.com/npavlov/go-password-manager/internal/server/service/meta"
	"github.com/npavlov/go-password-manager/internal/server/service/utils"
	testutils "github.com/npavlov/go-password-manager/internal/test_utils"
)

func setupMetadataService(t *testing.T) (*meta.Service, *testutils.MockDBStorage, context.Context) {
	t.Helper()

	logger := zerolog.New(nil)
	masterKey, _ := utils.GenerateRandomKey()
	storage := testutils.SetupMockUserStorage(masterKey)
	//nolint:exhaustruct
	cfg := &config.Config{}
	// Create test user
	userID := uuid.New()
	testUser := db.User{
		ID:       pgtype.UUID{Bytes: userID, Valid: true},
		Username: "testuser",
		Password: "hashed-password",
	}
	storage.AddTestUser(testUser)
	// Inject user ID into context
	ctx := testutils.InjectUserToContext(t.Context(), userID.String())

	return meta.NewMetadataService(&logger, storage, cfg), storage, ctx
}

func TestAddMetaInfo_Success(t *testing.T) {
	t.Parallel()

	svc, mockStorage, ctx := setupMetadataService(t)

	userID := testutils.GetUserIDFromContext(ctx)
	userIDPG := pgtype.UUID{
		Bytes: uuid.MustParse(userID),
		Valid: true,
	}

	newItem, err := mockStorage.StoreItem(ctx, userIDPG, db.ItemTypeCard)
	require.NoError(t, err)

	req := &pb.AddMetaInfoV1Request{
		ItemId: newItem.ID.String(),
		Metadata: map[string]string{
			"category": "finance",
			"priority": "high",
		},
	}

	resp, err := svc.AddMetaInfoV1(t.Context(), req)
	require.NoError(t, err)
	require.True(t, resp.GetSuccess())

	// Verify metadata was stored
	metaInfo, err := mockStorage.GetMetaInfo(t.Context(), newItem.ID.String())
	require.NoError(t, err)
	require.Len(t, metaInfo, 2)
	require.Equal(t, "finance", metaInfo[0].Value)
	require.Equal(t, "high", metaInfo[1].Value)
}

func TestAddMetaInfo_InvalidItem(t *testing.T) {
	t.Parallel()

	svc, _, ctx := setupMetadataService(t)

	req := &pb.AddMetaInfoV1Request{
		ItemId: "invalid-item-id",
		Metadata: map[string]string{
			"key": "value",
		},
	}

	_, err := svc.AddMetaInfoV1(ctx, req)
	require.Error(t, err)
	require.Contains(t, err.Error(), "error validating input")
}

func TestRemoveMetaInfo_Success(t *testing.T) {
	t.Parallel()

	svc, mockStorage, ctx := setupMetadataService(t)

	// Create test item with metadata
	itemID := uuid.NewString()
	_, err := mockStorage.AddMeta(ctx, itemID, "category", "finance")
	require.NoError(t, err)

	req := &pb.RemoveMetaInfoV1Request{
		ItemId: itemID,
		Key:    "category",
	}

	resp, err := svc.RemoveMetaInfoV1(ctx, req)
	require.NoError(t, err)
	require.True(t, resp.GetSuccess())

	// Verify metadata was removed
	metaInfo, err := mockStorage.GetMetaInfo(ctx, itemID)
	require.NoError(t, err)
	require.Empty(t, metaInfo)
}

func TestRemoveMetaInfo_NotFound(t *testing.T) {
	t.Parallel()

	svc, _, ctx := setupMetadataService(t)

	req := &pb.RemoveMetaInfoV1Request{
		ItemId: uuid.New().String(),
		Key:    "category",
	}

	_, err := svc.RemoveMetaInfoV1(ctx, req)
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to remove meta info")
}

func TestGetMetaInfo_Success(t *testing.T) {
	t.Parallel()

	svc, mockStorage, ctx := setupMetadataService(t)

	// Create test item with metadata
	itemID := uuid.NewString()
	_, err := mockStorage.AddMeta(ctx, itemID, "category", "finance")
	require.NoError(t, err)
	_, err = mockStorage.AddMeta(ctx, itemID, "priority", "high")
	require.NoError(t, err)

	req := &pb.GetMetaInfoV1Request{
		ItemId: itemID,
	}

	resp, err := svc.GetMetaInfoV1(t.Context(), req)
	require.NoError(t, err)
	require.Len(t, resp.GetMetadata(), 2)
	require.Equal(t, "finance", resp.GetMetadata()["category"])
	require.Equal(t, "high", resp.GetMetadata()["priority"])
}

func TestGetMetaInfo_Empty(t *testing.T) {
	t.Parallel()

	svc, mockStorage, ctx := setupMetadataService(t)

	userID := testutils.GetUserIDFromContext(ctx)
	userIDPG := pgtype.UUID{
		Bytes: uuid.MustParse(userID),
		Valid: true,
	}

	newItem, err := mockStorage.StoreItem(ctx, userIDPG, db.ItemTypeCard)
	require.NoError(t, err)

	req := &pb.GetMetaInfoV1Request{
		ItemId: newItem.ID.String(),
	}

	resp, err := svc.GetMetaInfoV1(ctx, req)
	require.NoError(t, err)
	require.Empty(t, resp.GetMetadata())
}

func TestAddMetaInfo_ValidationError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		req         *pb.AddMetaInfoV1Request
		expectedErr string
	}{
		{
			name: "empty item id",
			req: &pb.AddMetaInfoV1Request{
				ItemId:   "",
				Metadata: map[string]string{"key": "value"},
			},
			expectedErr: "error validating input",
		},
		{
			name: "empty metadata",
			req: &pb.AddMetaInfoV1Request{
				ItemId:   uuid.NewString(),
				Metadata: map[string]string{},
			},
			expectedErr: "error validating input",
		},
		{
			name: "invalid key",
			req: &pb.AddMetaInfoV1Request{
				ItemId:   uuid.NewString(),
				Metadata: map[string]string{"": "value"},
			},
			expectedErr: "error validating input",
		},
	}

	svc, _, ctx := setupMetadataService(t)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			_, err := svc.AddMetaInfoV1(ctx, tc.req)
			require.Error(t, err)
			require.Contains(t, err.Error(), tc.expectedErr)
		})
	}
}

func TestAddMetaInfo_StorageError(t *testing.T) {
	t.Parallel()

	svc, mockStorage, ctx := setupMetadataService(t)

	// Simulate storage failure
	mockStorage.CallError = errors.New("storage error")

	req := &pb.AddMetaInfoV1Request{
		ItemId:   uuid.NewString(),
		Metadata: map[string]string{"key": "value"},
	}

	_, err := svc.AddMetaInfoV1(ctx, req)
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to add meta info")
}

func TestRemoveMetaInfo_ValidationError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		req         *pb.RemoveMetaInfoV1Request
		expectedErr string
	}{
		{
			name: "empty item id",
			req: &pb.RemoveMetaInfoV1Request{
				ItemId: "",
				Key:    "key",
			},
			expectedErr: "error validating input",
		},
		{
			name: "empty key",
			req: &pb.RemoveMetaInfoV1Request{
				ItemId: uuid.NewString(),
				Key:    "",
			},
			expectedErr: "error validating input",
		},
	}

	svc, _, ctx := setupMetadataService(t)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			_, err := svc.RemoveMetaInfoV1(ctx, tc.req)
			require.Error(t, err)
			require.Contains(t, err.Error(), tc.expectedErr)
		})
	}
}

func TestRemoveMetaInfo_StorageError(t *testing.T) {
	t.Parallel()

	svc, mockStorage, ctx := setupMetadataService(t)

	// Simulate storage failure
	mockStorage.CallError = errors.New("database error")

	req := &pb.RemoveMetaInfoV1Request{
		ItemId: uuid.NewString(),
		Key:    "key",
	}

	_, err := svc.RemoveMetaInfoV1(ctx, req)
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to remove meta info")
}

func TestGetMetaInfo_ValidationError(t *testing.T) {
	t.Parallel()

	svc, _, ctx := setupMetadataService(t)

	req := &pb.GetMetaInfoV1Request{
		ItemId: "",
	}

	_, err := svc.GetMetaInfoV1(ctx, req)
	require.Error(t, err)
	require.Contains(t, err.Error(), "error validating input")
}

func TestGetMetaInfo_StorageError(t *testing.T) {
	t.Parallel()

	svc, mockStorage, ctx := setupMetadataService(t)

	mockStorage.CallError = errors.New("database error")

	req := &pb.GetMetaInfoV1Request{
		ItemId: uuid.NewString(),
	}

	_, err := svc.GetMetaInfoV1(ctx, req)
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to fetch metadata")
}
