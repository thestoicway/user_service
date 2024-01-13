package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/thestoicway/backend/user_service/internal/user/jsonwebtoken"
	"github.com/thestoicway/backend/user_service/internal/user/model"
)

// Refresh is a method of userService that takes a context and a refresh token,
// and returns a new pair of access and refresh tokens.
func (s *userServiceImpl) Refresh(ctx context.Context, refreshToken string) (tokenPair *jsonwebtoken.TokenPair, err error) {
	// Decode the refresh token to get the claims
	claims, err := s.JwtManager.DecodeToken(refreshToken)

	if err != nil {
		// If there's an error decoding the token, return the error
		return nil, err
	}

	// Parse the user ID from the claims
	uuid, err := uuid.Parse(claims.UserID)

	if err != nil {
		// If there's an error parsing the user ID, return the error
		return nil, err
	}

	// Get the current session using the ID from the claims
	currentSession, err := s.Session.GetSession(ctx, claims.ID)

	if err != nil {
		// If there's an error getting the session, return the error
		return nil, err
	}

	// Generate a new pair of access and refresh tokens
	pair, info, err := s.JwtManager.GenerateTokenPair(uuid)

	if err != nil {
		// If there's an error generating the token pair, return the error
		return nil, err
	}

	// Replace the old session with a new one that includes the new refresh token
	err = s.Session.ReplaceSession(ctx,
		currentSession,
		&model.Session{
			JwtID:          info.RefreshTokenID,
			RefreshToken:   pair.RefreshToken,
			ExpirationTime: info.RefreshExpirationTime,
		})

	if err != nil {
		// If there's an error replacing the session, return the error
		return nil, err
	}

	// Return the new pair of tokens
	return pair, nil
}
