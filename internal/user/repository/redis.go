package repository

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

type Redis struct {
	client *redis.Client
}

func NewRedis(client *redis.Client) *Redis {
	return &Redis{client: client}
}

func (r *Redis) StoreRefreshToken(ctx context.Context, token, userID string, expiresAt int64) error {
	key := fmt.Sprintf("refresh_token_%s", token)
	expiration := time.Unix(expiresAt, 0).Sub(time.Now())
	return r.client.Set(ctx, key, userID, expiration).Err()
}

func (r *Redis) GetUserIDByRefreshToken(ctx context.Context, token string) (string, error) {
	key := fmt.Sprintf("refresh_token_%s", token)
	return r.client.Get(ctx, key).Result()
}

func (r *Redis) DeleteRefreshToken(ctx context.Context, token string) error {
	key := fmt.Sprintf("refresh_token_%s", token)
	return r.client.Del(ctx, key).Err()
}
