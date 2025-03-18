package service

import (
	"context"
	"github.com/ashkan-maleki/kiwi-eats/internal/user/pb"
	"github.com/ashkan-maleki/kiwi-eats/internal/user/repository"
	"github.com/ashkan-maleki/kiwi-eats/internal/user/repository/entity"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

// Auth implements the AuthServiceServer interface
type Auth struct {
	userRepository *repository.User
	JWTSecret      string
}

// NewAuth creates a new Auth instance
func NewAuth(userRepository *repository.User, JWTSecret string) *Auth {
	return &Auth{userRepository: userRepository, JWTSecret: JWTSecret}
}

// Register handles user registration
func (s *Auth) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	// Check if user already exists
	existingUser, _ := s.userRepository.GetUserByEmail(ctx, req.Email)
	if existingUser != nil {
		return nil, status.Errorf(codes.AlreadyExists, "User with email %s already exists", req.Email)
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to hash password: %v", err)
	}
	// Create user entity
	user := &entity.User{
		Email:     req.Email,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
	}

	// Save to database
	userID, err := s.userRepository.CreateUser(ctx, user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to create user: %v", err)
	}

	// Return response
	return &pb.RegisterResponse{UserId: userID}, nil
}

func (s *Auth) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	// TODO: Implement login logic
	return nil, nil
}

func (s *Auth) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	// TODO: Implement token validation logic
	return nil, nil
}

func (s *Auth) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	// TODO: Implement token refresh logic
	return nil, nil
}

func (s *Auth) Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	// TODO: Implement logout logic
	return nil, nil
}
