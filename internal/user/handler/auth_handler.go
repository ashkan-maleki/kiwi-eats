package handler

import (
	"context"
	"github.com/ashkan-maleki/kiwi-eats/internal/user/pb"
	"github.com/ashkan-maleki/kiwi-eats/internal/user/service"
)

type Auth struct {
	authService *service.Auth
	pb.UnimplementedAuthServiceServer
}

func NewAuth(authService *service.Auth) *Auth {
	return &Auth{authService: authService}
}

func (h Auth) Register(ctx context.Context, request *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	return h.authService.Register(ctx, request)
}
