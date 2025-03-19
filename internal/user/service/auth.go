package service

import (
	"context"
	"github.com/ashkan-maleki/kiwi-eats/internal/user/pb"
	"github.com/ashkan-maleki/kiwi-eats/internal/user/repository/entity"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type (
	// UserRepo defines methods for user-related database operations.
	UserRepo interface {
		GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
		CreateUser(ctx context.Context, user *entity.User) (string, error)
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
	existingUser, _ := s.userRepo.GetUserByEmail(ctx, req.Email)
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
	userID, err := s.userRepo.CreateUser(ctx, user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to create user: %v", err)
	}

	// Return response
	return &pb.RegisterResponse{UserId: userID}, nil
}

func GenerateJWT(userID string, secretKey string, expiresAt int64) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     expiresAt, // 24-hour expiration
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

func (s *Auth) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	// Check if the user exists
	existingUser, _ := s.userRepo.GetUserByEmail(ctx, req.Email)
	if existingUser == nil {
		return nil, status.Errorf(codes.NotFound, "User with email %s does not exists", req.Email)
	}

	// Compare the stored hashed password with the provided password
	err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(req.Password))
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "Provided password does not match")
	}

	// If the password matches, generate a JWT access token, living 15 minutes
	accessTokenExpiresAt := time.Now().Add(time.Minute * 15).Unix()
	accessToken, err := GenerateJWT(existingUser.ID, s.JWTSecret, accessTokenExpiresAt)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to generate JWT: %v", err)
	}

	// Generate a refresh token with a 15-day expiration
	refreshTokenExpiresAt := time.Now().Add(time.Hour * 24 * 15).Unix() // 7 days
	refreshToken, err := GenerateJWT(existingUser.ID, s.JWTSecret, refreshTokenExpiresAt)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to generate refresh token: %v", err)
	}

	// Store the refresh token in Redis
	err = s.redisRepo.StoreRefreshToken(ctx, refreshToken, existingUser.ID, refreshTokenExpiresAt)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to store refresh token: %v", err)
	}

	// Store token metadata in PostgreSQL
	ipAddress := "192.168.1.1" // TODO: Extract from request context
	userAgent := "Mozilla/5.0" // TODO: Extract from request context
	err = s.tokenRepo.StoreTokenMetadata(ctx, existingUser.ID, refreshToken, ipAddress,
		userAgent, refreshTokenExpiresAt)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to store token metadata: %v", err)
	}
	// Return the access token
	return &pb.LoginResponse{
		JwtToken:     accessToken,
		RefreshToken: refreshToken,
		UserId:       existingUser.ID,
		Name:         existingUser.Email, // Todo: change email or name
		ExpiresAt:    accessTokenExpiresAt,
	}, nil
}

func ValidateJWT(tokenString string, secretKey string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
}

func (s *Auth) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	//Extract the JWT token from the request
	//Verify the token signature using the JWT secret key
	//Check if the token is expired
	//Extract user ID from token claims
	//Return success if valid, or an error if the token is invalid

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
