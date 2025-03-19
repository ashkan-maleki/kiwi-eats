package main

import (
	"fmt"
	"github.com/ashkan-maleki/kiwi-eats/internal/config"
	"github.com/ashkan-maleki/kiwi-eats/internal/user/handler"
	"github.com/ashkan-maleki/kiwi-eats/internal/user/pb"
	"github.com/ashkan-maleki/kiwi-eats/internal/user/repository"
	"github.com/ashkan-maleki/kiwi-eats/internal/user/service"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Could not load config: %v", err)
	}
	fmt.Println("App Name:", cfg.AppName)
	fmt.Println("Running on Port:", cfg.JWTSecret)
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	userRepo := repository.NewUser()
	redisRepo := repository.NewRedis(nil) // TODO: add redis client
	tokenRepo := repository.NewToken(nil) // TODO: add db client

	authService := service.NewAuth(userRepo, redisRepo, tokenRepo, cfg.JWTSecret)
	authHandler := handler.NewAuth(authService)

	grpcServer := grpc.NewServer()
	pb.RegisterAuthServiceServer(grpcServer, authHandler)

	log.Println("User service is running on port 50051...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
