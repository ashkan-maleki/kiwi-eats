package service

import (
	"context"
	"github.com/ashkan-maleki/kiwi-eats/internal/user/pb"
	"github.com/ashkan-maleki/kiwi-eats/internal/user/repository/entity"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthService_Login(t *testing.T) {
	password := "123456"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	require.NoError(t, err)

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

	tests := []struct {
		name        string
		email       string
		password    string
		mockUser    *entity.User
		mockError   error
		expectError bool
	}{
		{
			name:        "User not found",
			email:       "hacker@example.com",
			password:    "password",
			mockUser:    nil,
			mockError:   status.Errorf(codes.NotFound, "User not found"),
			expectError: true,
		},
		{
			name:        "Password does not match",
			email:       "test@example.com",
			password:    "wrong_password",
			mockUser:    &entity.User{ID: "123", Email: "test@example.com", Password: "hashed_password"},
			mockError:   nil,
			expectError: true,
		},
		{
			name:        "Login successful",
			email:       "test@example.com",
			password:    password,
			mockUser:    &entity.User{ID: "123", Email: "test@example.com", Password: string(hashedPassword)},
			mockError:   nil,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserRepo := &mockUserRepo{
				GetUserByEmailFunc: func(ctx context.Context, email string) (*entity.User, error) {
					return tt.mockUser, tt.mockError
				},
			}

			authService := NewAuth(mockUserRepo, mockRedisRepo, mockTokenRepo, "secret")

			resp, err := authService.Login(context.Background(), &pb.LoginRequest{
				Email:    tt.email,
				Password: tt.password,
			})

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotZero(t, resp.JwtToken)
				assert.NotZero(t, resp.RefreshToken)
			}
		})
	}
}
func TestAuthService_Register(t *testing.T) {
	tests := []struct {
		name             string
		userName         string
		email            string
		password         string
		mockUser         *entity.User
		mockError        error
		registeredUserID string
		expectError      bool
	}{
		{
			name:             "User already exists",
			userName:         "tester",
			email:            "test@example.com",
			password:         "password",
			registeredUserID: "",
			mockUser:         &entity.User{ID: "123", Email: "test@example.com", Password: "hashed_password"},
			mockError:        nil,
			expectError:      true,
		},
		{
			name:             "Register succeed",
			userName:         "tester",
			email:            "test@example.com",
			password:         "password",
			registeredUserID: "234",
			mockUser:         nil,
			mockError:        nil,
			expectError:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserRepo := &mockUserRepo{
				GetUserByEmailFunc: func(ctx context.Context, email string) (*entity.User, error) {
					return tt.mockUser, tt.mockError
				},
				CreateUserFunc: func(ctx context.Context, user *entity.User) (string, error) {
					return tt.registeredUserID, nil
				},
			}

			authService := NewAuth(mockUserRepo, &mockRedisRepo{}, &mockTokenRepo{}, "secret")

			resp, err := authService.Register(context.Background(), &pb.RegisterRequest{
				Name:     tt.userName,
				Email:    tt.email,
				Password: tt.password,
			})

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.registeredUserID, resp.UserId)
			}
		})
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
