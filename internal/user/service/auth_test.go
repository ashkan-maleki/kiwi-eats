package service

import (
	"context"
	"errors"
	"github.com/ashkan-maleki/kiwi-eats/internal/user/pb"
	"github.com/ashkan-maleki/kiwi-eats/internal/user/repository/entity"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthService_Login_UserNotFound(t *testing.T) {
	// Create mock repositories
	mockUserRepo := &mockUserRepo{
		GetUserByEmailFunc: func(ctx context.Context, email string) (*entity.User, error) {
			testEmail := "test@example.com"
			if email == email {
				return &entity.User{ID: "123", Email: testEmail, Password: "hashed_password"}, nil
			}
			return nil, status.Errorf(codes.NotFound, "User with email %s does not exists", testEmail)
		},
	}

	mockRedisRepo := &mockRedisRepo{
		StoreRefreshTokenFunc: func(ctx context.Context, token, userID string, expiresAt int64) error {
			return nil
		},
	}

	mockTokenRepo := &mockTokenRepo{
		StoreTokenMetadataFunc: func(ctx context.Context, userID, token, ipAddress, userAgent string, expiresAt int64) error {
			return nil
		},
	}

	// Create the service
	authService := NewAuth(mockUserRepo, mockRedisRepo, mockTokenRepo, "secret")

	// Test the Login method
	_, err := authService.Login(context.Background(), &pb.LoginRequest{
		Email:    "hacker@example.com",
		Password: "password",
	})
	assert.Error(t, err)
}

func TestAuthService_Login_PasswordDoesNotMatch(t *testing.T) {
	mockUserRepo := &mockUserRepo{
		GetUserByEmailFunc: func(ctx context.Context, email string) (*entity.User, error) {
			if email == "test@example.com" {
				return &entity.User{ID: "123", Email: "test@example.com", Password: "hashed_password"}, nil
			}
			return nil, errors.New("user not found")
		},
	}

	mockRedisRepo := &mockRedisRepo{
		StoreRefreshTokenFunc: func(ctx context.Context, token, userID string, expiresAt int64) error {
			return nil
		},
	}

	mockTokenRepo := &mockTokenRepo{
		StoreTokenMetadataFunc: func(ctx context.Context, userID, token, ipAddress, userAgent string, expiresAt int64) error {
			return nil
		},
	}

	// Create the service
	authService := NewAuth(mockUserRepo, mockRedisRepo, mockTokenRepo, "secret")

	// Test the Login method
	_, err := authService.Login(context.Background(), &pb.LoginRequest{
		Email:    "test@example.com",
		Password: "password",
	})
	assert.Error(t, err)
}

func TestAuthService_Login_PasswordDoesNotMatch1(t *testing.T) {
	mockUserRepo := &mockUserRepo{
		GetUserByEmailFunc: func(ctx context.Context, email string) (*entity.User, error) {
			if email == "test@example.com" {
				return &entity.User{ID: "123", Email: "test@example.com", Password: "hashed_password"}, nil
			}
			return nil, errors.New("user not found")
		},
	}

	mockRedisRepo := &mockRedisRepo{
		StoreRefreshTokenFunc: func(ctx context.Context, token, userID string, expiresAt int64) error {
			return nil
		},
	}

	mockTokenRepo := &mockTokenRepo{
		StoreTokenMetadataFunc: func(ctx context.Context, userID, token, ipAddress, userAgent string, expiresAt int64) error {
			return nil
		},
	}

	// Create the service
	authService := NewAuth(mockUserRepo, mockRedisRepo, mockTokenRepo, "secret")

	// Test the Login method
	resp, err := authService.Login(context.Background(), &pb.LoginRequest{
		Email:    "test@example.com",
		Password: "password",
	})
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}
	if resp.UserId != "123" {
		t.Errorf("Expected user ID 123, got %s", resp.UserId)
	}
}

type mockUserRepo struct {
	GetUserByEmailFunc func(ctx context.Context, email string) (*entity.User, error)
	CreateUserFunc     func(ctx context.Context, user *entity.User) (string, error)
}

func (m *mockUserRepo) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	return m.GetUserByEmailFunc(ctx, email)
}

func (m *mockUserRepo) CreateUser(ctx context.Context, user *entity.User) (string, error) {
	return m.CreateUserFunc(ctx, user)
}

type mockRedisRepo struct {
	StoreRefreshTokenFunc func(ctx context.Context, token, userID string, expiresAt int64) error
}

func (m *mockRedisRepo) StoreRefreshToken(ctx context.Context, token, userID string, expiresAt int64) error {
	return m.StoreRefreshTokenFunc(ctx, token, userID, expiresAt)
}

type mockTokenRepo struct {
	StoreTokenMetadataFunc func(ctx context.Context, userID, token, ipAddress, userAgent string, expiresAt int64) error
}

func (m *mockTokenRepo) StoreTokenMetadata(ctx context.Context, userID, token, ipAddress, userAgent string, expiresAt int64) error {
	return m.StoreTokenMetadataFunc(ctx, userID, token, ipAddress, userAgent, expiresAt)
}
