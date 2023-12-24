package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/thestoicway/backend/user_service/internal/jsonwebtoken"
)

func (svc *userService) Refresh(ctx context.Context, refreshToken string) (tokenPair *jsonwebtoken.TokenPair, err error) {
	claims, err := svc.JwtManager.DecodeToken(refreshToken)

	if err != nil {
		return nil, err
	}

	uuid, err := uuid.Parse(claims.UserID)

	if err != nil {
		return nil, err
	}

	pair, err := svc.JwtManager.GenerateTokenPair(uuid)

	if err != nil {
		return nil, err
	}

	return pair, nil
}
