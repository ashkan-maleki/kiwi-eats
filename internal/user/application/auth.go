package domain

import (
	"context"
	"github.com/ashkan-maleki/kiwi-eats/internal/user/pb"
	"github.com/ashkan-maleki/kiwi-eats/internal/user/repository/entity"
	"github.com/ashkan-maleki/kiwi-eats/pkg/auth"
	"github.com/ashkan-maleki/kiwi-eats/pkg/grpc"
	"github.com/ashkan-maleki/kiwi-eats/pkg/logger"
	"go.uber.org/zap"

	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type (

	// Auth implements the AuthServiceServer interface
	Auth struct {
		userRepo  UserRepo
		redisRepo RedisRepo
		tokenRepo TokenRepo
		JWTSecret string
	}
)

// NewAuth creates a new Auth instance
func NewAuth(userRepo UserRepo, redisRepo RedisRepo, tokenRepo TokenRepo, JWTSecret string) *Auth {
	return &Auth{
		userRepo:  userRepo,
		redisRepo: redisRepo,
		tokenRepo: tokenRepo,
		JWTSecret: JWTSecret,
	}
}

// Register handles user registration
func (s *Auth) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	// Check if user already exists
	existingUser, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot find user by email: %v", err)
	}
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
		Name:      req.Name,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
	}

	// Save to database
	userID, err := s.userRepo.CreateUser(ctx, user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to create user: %v", err)
	}

	// Return response
	return &pb.RegisterResponse{UserId: userID}, nil
}

func (s *Auth) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	// Check if the user exists
	existingUser, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil || existingUser == nil {
		logger.Logger2.Error("User not found", zap.String("email", req.Email), zap.Error(err))
		return nil, status.Errorf(codes.NotFound, "User with email %s does not exist", req.Email)
	}

	// Compare the stored hashed password with the provided password
	err = bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(req.Password))
	if err != nil {
		logger.Logger2.Error("Password does not match", zap.String("email", req.Email))
		return nil, status.Errorf(codes.Unauthenticated, "Provided password does not match")
	}

	accessTokenExpiresAt := time.Now().Add(time.Minute * 15).Unix()     // 15 minutes
	refreshTokenExpiresAt := time.Now().Add(time.Hour * 24 * 15).Unix() // 15 days

	accessToken, refreshToken, err := auth.GenerateJWT(existingUser.ID, s.JWTSecret, accessTokenExpiresAt, refreshTokenExpiresAt)
	if err != nil {
		logger.Logger2.Error("Failed to generate JWT tokens", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "Failed to generate JWT: %v", err)
	}

	// Store the refresh token in Redis
	err = s.redisRepo.StoreRefreshToken(ctx, refreshToken, existingUser.ID, refreshTokenExpiresAt)
	if err != nil {
		logger.Logger2.Error("Failed to store refresh token in Redis", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "Failed to store refresh token: %v", err)
	}

	// Extract IP address and user agent
	md, err := grpc.Metadata(ctx, grpc.UserAgent, grpc.IP)
	if err != nil {
		logger.Logger2.Error("Failed to extract metadata", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "Failed to get metadata: %v", err)
	}

	// Store token metadata in PostgreSQL
	err = s.tokenRepo.StoreTokenMetadata(ctx, existingUser.ID, refreshToken, md.IpAddress(), md.UserAgent(), refreshTokenExpiresAt)
	if err != nil {
		logger.Logger2.Error("Failed to store token metadata in PostgreSQL", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "Failed to store token metadata: %v", err)
	}

	logger.Logger2.Info("User logged in successfully", zap.String("user_id", existingUser.ID))
	return &pb.LoginResponse{
		JwtToken:     accessToken,
		RefreshToken: refreshToken,
		UserId:       existingUser.ID,
		Name:         existingUser.Name,
		ExpiresAt:    accessTokenExpiresAt,
	}, nil
}

func (s *Auth) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	if err := auth.ValidateJWTToken(ctx, req.Token, s.JWTSecret); err != nil {
		logger.Logger2.Error("Failed to validate JWT token", zap.Error(err))
		return nil, status.Errorf(codes.Unauthenticated, "failed to authorize: %v", err)
	}
	return &pb.ValidateTokenResponse{IsValid: true}, nil
}

func (s *Auth) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	// TODO: Implement token refresh logic
	return nil, nil
}

func (s *Auth) Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	// TODO: Implement logout logic
	return nil, nil
}
