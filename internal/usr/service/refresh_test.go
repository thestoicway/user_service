package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/thestoicway/user_service/internal/usr/jsonwebtoken"
	"github.com/thestoicway/user_service/internal/usr/model"
	"github.com/thestoicway/user_service/internal/usr/service"
	sessionstorage "github.com/thestoicway/user_service/internal/usr/session_storage"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

func TestRefresh(t *testing.T) {
	t.Run("refreshes token pair", func(t *testing.T) {
		t.Parallel()

		jwtManager := jsonwebtoken.NewJwtManager(zaptest.NewLogger(t).Sugar(), "secret")

		oldJwt, info, _ := jwtManager.GenerateTokenPair(uuid.New())

		oldSession := &model.Session{
			JwtID:          info.RefreshTokenID,
			RefreshToken:   oldJwt.RefreshToken,
			ExpirationTime: time.Duration(1) * time.Hour,
		}

		// Create a new session storage with the old session.
		sStorage := sessionstorage.NewInMemoryDatabase(&sessionstorage.InMemoryDatabaseParams{
			Sessions: map[string]*model.Session{
				oldSession.JwtID: oldSession,
			},
		})

		// Create a new user service with the session storage.
		userService := service.NewUserService(&service.UserServiceParams{
			Logger:     zap.NewNop().Sugar(),
			JwtManager: jwtManager,
			// Database is not needed for this test.
			Database: nil,
			Session:  sStorage,
		})

		// Refresh the token pair.
		tokenPair, err := userService.Refresh(context.Background(), oldSession.RefreshToken)

		if assert.NoError(t, err) {
			// Check that the new refresh token is the same as the one in the new session.
			assert.NotEqual(t, tokenPair.RefreshToken, oldJwt.RefreshToken)
		}
	})

	t.Run("returns error when refresh token is invalid", func(t *testing.T) {
		t.Parallel()

		// Setup
		jwtManager := jsonwebtoken.NewJwtManager(zaptest.NewLogger(t).Sugar(), "secret")
		sStorage := sessionstorage.NewInMemoryDatabase(&sessionstorage.InMemoryDatabaseParams{})
		userService := service.NewUserService(&service.UserServiceParams{
			Logger:     zap.NewNop().Sugar(),
			JwtManager: jwtManager,
			Database:   nil,
			Session:    sStorage,
		})

		// Call Refresh with an invalid token
		_, err := userService.Refresh(context.Background(), "invalid_token")

		// Assert
		assert.Error(t, err)
	})
}
