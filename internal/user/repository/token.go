package repository

import (
	"context"
	"database/sql"
	"time"
)

type Token struct {
	db *sql.DB
}

func NewToken(db *sql.DB) *Token {
	return &Token{db: db}
}

//TODO: CREATE TABLE token_metadata (
//id SERIAL PRIMARY KEY,
//user_id TEXT NOT NULL,
//token TEXT NOT NULL UNIQUE,
//ip_address TEXT,
//user_agent TEXT,
//expires_at TIMESTAMP NOT NULL
//);

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
func (r *Token) StoreTokenMetadata(ctx context.Context, userID, token, ipAddress, userAgent string, expiresAt int64) error {
	query := `  
		INSERT INTO token_metadata (user_id, token, ip_address, user_agent, expires_at)  
		VALUES ($1, $2, $3, $4, $5)  
	`
	_, err := r.db.ExecContext(ctx, query, userID, token, ipAddress, userAgent, time.Unix(expiresAt, 0))
	return err
}

func (r *Token) GetTokenMetadata(ctx context.Context, token string) (userID, ipAddress, userAgent string, expiresAt int64, err error) {
	query := `  
		SELECT user_id, ip_address, user_agent, expires_at  
		FROM token_metadata  
		WHERE token = $1  
	`
	var expiresAtTime time.Time
	err = r.db.QueryRowContext(ctx, query, token).Scan(&userID, &ipAddress, &userAgent, &expiresAtTime)
	if err != nil {
		return "", "", "", 0, err
	}
	return userID, ipAddress, userAgent, expiresAtTime.Unix(), nil
}

func (r *Token) DeleteTokenMetadata(ctx context.Context, token string) error {
	query := `DELETE FROM token_metadata WHERE token = $1`
	_, err := r.db.ExecContext(ctx, query, token)
	return err
}
