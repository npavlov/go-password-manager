package item_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	pb "github.com/npavlov/go-password-manager/gen/proto/item"
	"github.com/npavlov/go-password-manager/internal/server/config"
	"github.com/npavlov/go-password-manager/internal/server/db"
	"github.com/npavlov/go-password-manager/internal/server/service/item"
	"github.com/npavlov/go-password-manager/internal/server/service/utils"
	testutils "github.com/npavlov/go-password-manager/internal/test_utils"
)

func setupItemService(t *testing.T) (*item.Service, *testutils.MockDBStorage, context.Context) {
	t.Helper()

	logger := zerolog.New(nil)
	masterKey, _ := utils.GenerateRandomKey()
	storage := testutils.SetupMockUserStorage(masterKey)
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

	return item.NewItemService(&logger, storage, cfg), storage, ctx
}

func TestGetItems_Success(t *testing.T) {
	t.Parallel()

	svc, storage, ctx := setupItemService(t)

	// Add test items
	userID := testutils.GetUserIDFromContext(ctx)
	now := time.Now()

	items := []db.Item{
		{
			ID:        pgtype.UUID{Bytes: uuid.New(), Valid: true},
			UserID:    pgtype.UUID{Bytes: uuid.MustParse(userID), Valid: true},
			Type:      db.ItemTypeCard,
			CreatedAt: pgtype.Timestamp{Time: now, Valid: true},
		},
		{
			ID:        pgtype.UUID{Bytes: uuid.New(), Valid: true},
			UserID:    pgtype.UUID{Bytes: uuid.MustParse(userID), Valid: true},
			Type:      db.ItemTypePassword,
			CreatedAt: pgtype.Timestamp{Time: now.Add(-time.Hour), Valid: true},
		},
	}

	// Store items directly in mock storage
	for _, item := range items {
		_, err := storage.StoreItem(t.Context(), item.UserID, item.Type)
		require.NoError(t, err)
	}

	// Test request
	req := &pb.GetItemsRequest{
		Page:     1,
		PageSize: 10,
	}

	resp, err := svc.GetItems(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, int32(2), resp.GetTotalCount())
}

func TestGetItems_Pagination(t *testing.T) {
	t.Parallel()

	svc, storage, ctx := setupItemService(t)

	// Add 15 test items
	userID := testutils.GetUserIDFromContext(ctx)
	userIDPG := pgtype.UUID{
		Bytes: uuid.MustParse(userID),
		Valid: true,
	}

	for range 15 {
		_, err := storage.StoreItem(t.Context(), userIDPG, db.ItemTypeText)
		require.NoError(t, err)
	}

	// First page
	req1 := &pb.GetItemsRequest{
		Page:     1,
		PageSize: 10,
	}
	resp1, err := svc.GetItems(ctx, req1)
	require.NoError(t, err)
	require.Equal(t, int32(10), resp1.GetTotalCount())
	require.Len(t, resp1.GetItems(), 10)

	// Second page
	req2 := &pb.GetItemsRequest{
		Page:     2,
		PageSize: 10,
	}
	resp2, err := svc.GetItems(ctx, req2)
	require.NoError(t, err)
	require.Equal(t, int32(5), resp2.GetTotalCount())
	require.Len(t, resp2.GetItems(), 5)
}

func TestGetItems_EmptyResult(t *testing.T) {
	t.Parallel()

	svc, _, ctx := setupItemService(t)

	req := &pb.GetItemsRequest{
		Page:     1,
		PageSize: 10,
	}

	resp, err := svc.GetItems(ctx, req)
	require.NoError(t, err)
	require.Equal(t, int32(0), resp.GetTotalCount())
	require.Empty(t, resp.GetItems())
}

func TestGetItems_InvalidRequest(t *testing.T) {
	t.Parallel()

	svc, _, ctx := setupItemService(t)

	// Page = 0 is invalid
	req := &pb.GetItemsRequest{
		Page:     0,
		PageSize: 10,
	}

	_, err := svc.GetItems(ctx, req)
	require.Error(t, err)
	require.Contains(t, err.Error(), "error validating input")
}

func TestGetItems_NoUserContext(t *testing.T) {
	t.Parallel()

	svc, _, _ := setupItemService(t)

	req := &pb.GetItemsRequest{
		Page:     1,
		PageSize: 10,
	}

	_, err := svc.GetItems(t.Context(), req)
	require.Error(t, err)
	require.Contains(t, err.Error(), "error getting user id")
}
