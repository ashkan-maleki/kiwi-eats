package handler

import (
	"context"
	"github.com/ashkan-maleki/kiwi-eats/internal/user/pb"
)

type Auth struct {
	authService *domain.Auth
	pb.UnimplementedAuthServiceServer
}

func NewAuth(authService *domain.Auth) *Auth {
	return &Auth{authService: authService}
}

func (h Auth) Register(ctx context.Context, request *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	return h.authService.Register(ctx, request)
}

func (h Auth) Login(ctx context.Context, request *pb.LoginRequest) (*pb.LoginResponse, error) {
	return h.authService.Login(ctx, request)
}

func (h Auth) ValidateToken(ctx context.Context, request *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	return h.authService.ValidateToken(ctx, request)
}

func (h Auth) RefreshToken(ctx context.Context, request *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	return h.authService.RefreshToken(ctx, request)
}

func (h Auth) Logout(ctx context.Context, request *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	return h.authService.Logout(ctx, request)
}
