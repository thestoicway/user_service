package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	customerrors "github.com/thestoicway/custom_errors"
	"github.com/thestoicway/user_service/internal/usr/jsonwebtoken"
	mocks "github.com/thestoicway/user_service/internal/usr/mocks"
	"github.com/thestoicway/user_service/internal/usr/model"
	"github.com/thestoicway/user_service/internal/usr/service"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap/zaptest"
)

func TestSignIn(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockDatabase := mocks.NewMockUserDatabase(mockCtrl)
	mockJwtManager := mocks.NewMockJwtManager(mockCtrl)
	mockSession := mocks.NewMockSessionDatabase(mockCtrl)

	userService := service.NewUserService(
		&service.UserServiceParams{
			Database:   mockDatabase,
			JwtManager: mockJwtManager,
			Session:    mockSession,
			Logger:     zaptest.NewLogger(t).Sugar(),
		},
	)

	ctx := context.Background()
	user := &model.User{Email: "test@example.com", Password: "password"}

	t.Run("error on GetUserByEmail", func(t *testing.T) {
		mockDatabase.EXPECT().GetUserByEmail(ctx, user.Email).Return(nil, errors.New("database error"))

		_, err := userService.SignIn(ctx, user)

		assert.Error(t, err)
	})

	t.Run("wrong password", func(t *testing.T) {
		// bcrypt hash of "hello"
		password := "$2a$12$0tfcCrafm/TctMHT69YkbOOeiWqUGN.uglp.hTUF8FfumtPdrfWJS"
		mockDatabase.EXPECT().GetUserByEmail(ctx, user.Email).Return(&model.UserDB{
			PasswordHash: password,
			Email:        user.Email,
		}, nil)

		_, err := userService.SignIn(ctx, user)

		if err, ok := err.(*customerrors.CustomError); ok {
			assert.Equal(t, err.Code, customerrors.ErrWrongCredentials)
		} else {
			assert.Fail(t, "error is not of type WrongCredentialsError")
		}
	})

	t.Run("error on GenerateTokenPair", func(t *testing.T) {
		mockDatabase.EXPECT().GetUserByEmail(ctx, user.Email).Return(&model.UserDB{
			// bcrypt hash of "password"
			PasswordHash: "$2y$10$THUQCeWjc8kIN2fYOz3ace43UVhfC1CaIqMiOU2.PxRqSTrbrXkE2",
			Email:        user.Email,
		}, nil)

		jwtError := errors.New("jwt error")

		mockJwtManager.EXPECT().GenerateTokenPair(gomock.Any()).Return(nil, nil, jwtError)

		_, err := userService.SignIn(ctx, user)

		assert.ErrorAs(t, err, &jwtError)
	})

	t.Run("error on AddSession", func(t *testing.T) {
		mockDatabase.EXPECT().GetUserByEmail(ctx, user.Email).Return(&model.UserDB{
			// bcrypt hash of "password"
			PasswordHash: "$2y$10$THUQCeWjc8kIN2fYOz3ace43UVhfC1CaIqMiOU2.PxRqSTrbrXkE2",
			Email:        user.Email,
		}, nil)
		mockJwtManager.EXPECT().GenerateTokenPair(gomock.Any()).Return(&jsonwebtoken.TokenPair{}, &jsonwebtoken.AdditionalInfo{}, nil)
		mockSession.EXPECT().AddSession(ctx, gomock.Any()).Return(errors.New("session error"))

		_, err := userService.SignIn(ctx, user)
		assert.Error(t, err)
	})

	t.Run("success", func(t *testing.T) {
		mockDatabase.EXPECT().GetUserByEmail(ctx, user.Email).Return(&model.UserDB{
			// bcrypt hash of "password"
			PasswordHash: "$2y$10$THUQCeWjc8kIN2fYOz3ace43UVhfC1CaIqMiOU2.PxRqSTrbrXkE2",
		}, nil)
		mockJwtManager.EXPECT().GenerateTokenPair(gomock.Any()).Return(&jsonwebtoken.TokenPair{}, &jsonwebtoken.AdditionalInfo{}, nil)
		mockSession.EXPECT().AddSession(ctx, gomock.Any()).Return(nil)

		_, err := userService.SignIn(ctx, user)
		assert.NoError(t, err)
	})
}
