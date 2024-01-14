package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/thestoicway/user_service/internal/usr/jsonwebtoken"
	mocks "github.com/thestoicway/user_service/internal/usr/mocks"
	"github.com/thestoicway/user_service/internal/usr/model"
	"github.com/thestoicway/user_service/internal/usr/service"
	"go.uber.org/mock/gomock"
)

func TestSignUp(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDatabase := mocks.NewMockUserDatabase(ctrl)
	mockJwtManager := mocks.NewMockJwtManager(ctrl)
	mockSession := mocks.NewMockSessionDatabase(ctrl)

	s := service.NewUserService(&service.UserServiceParams{
		Database:   mockDatabase,
		JwtManager: mockJwtManager,
		Session:    mockSession,
	})

	user := &model.User{
		Email:    "test@example.com",
		Password: "password",
	}

	ctx := context.Background()

	t.Run("successful sign up", func(t *testing.T) {
		mockDatabase.EXPECT().InsertUser(gomock.Any(), gomock.Any()).Return(&uuid.UUID{}, nil)
		mockJwtManager.EXPECT().GenerateTokenPair(gomock.Any()).Return(&jsonwebtoken.TokenPair{}, &jsonwebtoken.AdditionalInfo{}, nil)
		mockSession.EXPECT().AddSession(gomock.Any(), gomock.Any()).Return(nil)

		_, err := s.SignUp(ctx, user)
		assert.NoError(t, err)
	})

	t.Run("error when generating password hash", func(t *testing.T) {
		user.Password = string(make([]byte, 1001)) // bcrypt can handle up to 72 bytes

		_, err := s.SignUp(ctx, user)
		assert.Error(t, err)
	})

	t.Run("error when inserting user into the database", func(t *testing.T) {
		user.Password = "password"
		mockDatabase.EXPECT().InsertUser(gomock.Any(), gomock.Any()).Return(&uuid.UUID{}, errors.New("database error"))

		_, err := s.SignUp(ctx, user)
		assert.Error(t, err)
	})

	t.Run("error when generating token pair", func(t *testing.T) {
		mockDatabase.EXPECT().InsertUser(gomock.Any(), gomock.Any()).Return(&uuid.UUID{}, nil)
		mockJwtManager.EXPECT().GenerateTokenPair(gomock.Any()).Return(nil, nil, errors.New("token error"))

		_, err := s.SignUp(ctx, user)
		assert.Error(t, err)
	})

	t.Run("error when adding session", func(t *testing.T) {
		mockDatabase.EXPECT().InsertUser(gomock.Any(), gomock.Any()).Return(&uuid.UUID{}, nil)
		mockJwtManager.EXPECT().GenerateTokenPair(gomock.Any()).Return(&jsonwebtoken.TokenPair{}, &jsonwebtoken.AdditionalInfo{}, nil)
		mockSession.EXPECT().AddSession(gomock.Any(), gomock.Any()).Return(errors.New("session error"))

		_, err := s.SignUp(ctx, user)
		assert.Error(t, err)
	})

	t.Run("error when user already exists", func(t *testing.T) {
		mockDatabase.EXPECT().InsertUser(gomock.Any(), gomock.Any()).Return("", errors.New("user already exists"))

		_, err := s.SignUp(ctx, user)
		assert.Error(t, err)
	})

	t.Run("error when email is invalid", func(t *testing.T) {
		user.Email = "invalid email"

		mockDatabase.EXPECT().InsertUser(gomock.Any(), gomock.Any()).Return("", nil)
		mockJwtManager.EXPECT().GenerateTokenPair(gomock.Any()).Return(&jsonwebtoken.TokenPair{}, &jsonwebtoken.AdditionalInfo{}, nil)
		mockSession.EXPECT().AddSession(gomock.Any(), gomock.Any()).Return(nil)

		_, err := s.SignUp(ctx, user)
		assert.Error(t, err)
	})

	t.Run("error when password is too short", func(t *testing.T) {
		user.Password = "short"

		mockDatabase.EXPECT().InsertUser(gomock.Any(), gomock.Any()).Return("", nil)
		mockJwtManager.EXPECT().GenerateTokenPair(gomock.Any()).Return(&jsonwebtoken.TokenPair{}, &jsonwebtoken.AdditionalInfo{}, nil)
		mockSession.EXPECT().AddSession(gomock.Any(), gomock.Any()).Return(nil)

		_, err := s.SignUp(ctx, user)
		assert.Error(t, err)
	})
}
