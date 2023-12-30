package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/thestoicway/backend/user_service/internal/jsonwebtoken"
	"github.com/thestoicway/backend/user_service/internal/user/model"
)

func (s *userService) Refresh(ctx context.Context, refreshToken string) (tokenPair *jsonwebtoken.TokenPair, err error) {
	claims, err := s.JwtManager.DecodeToken(refreshToken)

	if err != nil {
		return nil, err
	}

	uuid, err := uuid.Parse(claims.UserID)

	if err != nil {
		return nil, err
	}

	currentSession, err := s.Session.GetSession(ctx, claims.ID)

	if err != nil {
		return nil, err
	}

	pair, info, err := s.JwtManager.GenerateTokenPair(uuid)

	if err != nil {
		return nil, err
	}

	// Delete the previous session of the user.
	// So that the refresh token can't be used anymore.
	err = s.Session.DeleteSession(ctx, currentSession.JwtID)

	if err != nil {
		return nil, err
	}

	// Creates a session in Redis where the key is ID of the
	// refresh token and the value is the refresh token itself.
	err = s.Session.AddSession(ctx, &model.Session{
		JwtID:          info.RefreshTokenID,
		RefreshToken:   pair.RefreshToken,
		ExpirationTime: info.RefreshExpirationTime,
	})

	if err != nil {
		return nil, err
	}

	return pair, nil
}
