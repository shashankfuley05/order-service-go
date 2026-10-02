package service

import (
	"fmt"

	"github.com/shashank/order-service/internal/auth"
	"github.com/shashank/order-service/internal/dto"
)

type AuthService interface {
	GenerateToken(tokenRequest *dto.TokenRequest) (*dto.TokenResponse, error)
}

type AuthServiceImpl struct {
}

func (s *AuthServiceImpl) GenerateToken(tokenRequest *dto.TokenRequest) (*dto.TokenResponse, error) {

	token, err := auth.GenerateTokenWithClaims(tokenRequest.UserID, "USER")

	if err != nil {
		return nil, fmt.Errorf("Token cannot be generated ; %w", err)
	}

	return &dto.TokenResponse{
		Token: token,
	}, nil
}
