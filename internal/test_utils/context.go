package testutils

import (
	"context"
)

// InjectUserToContext adds userID and encryptionKey into context — for tests or middleware.
func InjectUserToContext(ctx context.Context, userID string) context.Context {
	ctx = context.WithValue(ctx, "user_id", userID)

	return ctx
}

// GetUserIDFromContext retrieves the userID from the context.
// Returns the userID and a boolean indicating if it was found.
func GetUserIDFromContext(ctx context.Context) string {
	//nolint:forcetypeassert
	userID := ctx.Value("user_id").(string)

	return userID
}
