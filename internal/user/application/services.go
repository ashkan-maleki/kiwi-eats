package domain

import (
	"context"
	"github.com/ashkan-maleki/kiwi-eats/internal/user/domain"
)

type (
	// UserRepo defines methods for user-related database operations.
	UserRepo interface {
		GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
		CreateUser(ctx context.Context, user *domain.User) (string, error)
	}

	// RedisRepo defines methods for Redis operations.
	RedisRepo interface {
		StoreRefreshToken(ctx context.Context, token, userID string, expiresAt int64) error
	}

	// TokenRepo defines methods for token metadata operations.
	TokenRepo interface {
		// StoreTokenMetadata stores token metadata in the database.
		//
		// Parameters:
		//   - ctx: Context for request cancellation and timeouts.
		//   - userID: The ID of the user associated with the token.
		//   - token: The refresh token to store.
		//   - ipAddress: The IP address of the client making the request.
		//   - userAgent: The user agent of the client making the request.
		//   - expiresAt: The expiration time of the token (Unix timestamp).
		//
		// Returns:
		//   - error: An error if the operation fails.
		StoreTokenMetadata(ctx context.Context, userID, token, ipAddress, userAgent string, expiresAt int64) error
	}
)
