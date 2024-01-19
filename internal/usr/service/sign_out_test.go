package service_test

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/thestoicway/user_service/internal/usr/jsonwebtoken"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thestoicway/user_service/internal/usr/mocks"
	"github.com/thestoicway/user_service/internal/usr/service"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap/zaptest"
)

func TestSignOut(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockJwtManager := mocks.NewMockJwtManager(mockCtrl)
	mockSession := mocks.NewMockSessionDatabase(mockCtrl)

	userService := service.NewUserService(
		&service.UserServiceParams{
			JwtManager: mockJwtManager,
			Session:    mockSession,
			Logger:     zaptest.NewLogger(t).Sugar(),
		},
	)

	ctx := context.Background()
	refreshToken := "some_refresh_token"

	t.Run("error on DecodeToken", func(t *testing.T) {
		jwtError := errors.New("jwt error")
		mockJwtManager.EXPECT().DecodeToken(refreshToken).Return(nil, jwtError)

		err := userService.SignOut(ctx, refreshToken)

		assert.ErrorAs(t, err, &jwtError)
	})

	t.Run("error on DeleteSession", func(t *testing.T) {
		sessionError := errors.New("session error")
		mockJwtManager.EXPECT().DecodeToken(refreshToken).Return(&jsonwebtoken.CustomClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				ID: "123",
			},
		}, nil)
		mockSession.EXPECT().DeleteSession(ctx, "123").Return(sessionError)

		err := userService.SignOut(ctx, refreshToken)

		assert.ErrorAs(t, err, &sessionError)
	})

	t.Run("success", func(t *testing.T) {
		mockJwtManager.EXPECT().DecodeToken(refreshToken).Return(&jsonwebtoken.CustomClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				ID: "123",
			},
		}, nil)
		mockSession.EXPECT().DeleteSession(ctx, "123").Return(nil)

		err := userService.SignOut(ctx, refreshToken)

		assert.NoError(t, err)
	})
}
