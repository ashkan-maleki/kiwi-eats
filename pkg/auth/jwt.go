package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/ashkan-maleki/kiwi-eats/internal/user/repository/entity"
	"github.com/ashkan-maleki/kiwi-eats/pkg/logger"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"time"
)

func ValidateJWTToken(_ context.Context, tokenString, secretKey string) error {
	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		logger.Logger.Error("Failed to parse JWT token", zap.Error(err))
		return fmt.Errorf("failed to parse token: %v", err)
	}

	// Validate the token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		logger.Logger.Error("Invalid JWT token")
		return errors.New("invalid token")
	}

	// Check expiration
	exp, ok := claims["exp"].(float64)
	if !ok {
		logger.Logger.Error("Invalid expiration claim in JWT token")
		return errors.New("invalid expiration claim")
	}
	if time.Now().Unix() > int64(exp) {
		logger.Logger.Error("JWT token expired")
		return errors.New("token expired")
	}

	return nil
}

func generate(userID, secretKey string, expiresAt int64) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     expiresAt,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

func GenerateJWT(existingUser *entity.User, secretKey string, accessTokenExpiresAt, refreshTokenExpiresAt int64) (string, string, error) {
	// Generate a JWT access token
	accessToken, err := generate(existingUser.ID, secretKey, accessTokenExpiresAt)
	if err != nil {
		logger.Logger.Error("Failed to generate access token", zap.Error(err))
		return "", "", fmt.Errorf("failed to generate access token: %w", err)
	}

	// Generate a refresh token
	refreshToken, err := generate(existingUser.ID, secretKey, refreshTokenExpiresAt)
	if err != nil {
		logger.Logger.Error("Failed to generate refresh token", zap.Error(err))
		return "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	logger.Logger.Info("Generated JWT tokens",
		zap.String("user_id", existingUser.ID),
		zap.Int64("access_token_expires_at", accessTokenExpiresAt),
		zap.Int64("refresh_token_expires_at", refreshTokenExpiresAt),
	)
	return accessToken, refreshToken, nil
}
