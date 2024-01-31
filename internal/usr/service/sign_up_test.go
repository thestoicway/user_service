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
	"go.uber.org/zap/zaptest"
)

func TestSignUp(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDatabase := mocks.NewMockUserDatabase(ctrl)
	mockJwtManager := mocks.NewMockJwtManager(ctrl)
	mockSession := mocks.NewMocksessionstorage(ctrl)

	s := service.NewUserService(&service.UserServiceParams{
		Database:   mockDatabase,
		JwtManager: mockJwtManager,
		Session:    mockSession,
		Logger:     zaptest.NewLogger(t).Sugar(),
	})

	ctx := context.Background()

	t.Run("successful sign up", func(t *testing.T) {
		mockDatabase.EXPECT().InsertUser(gomock.Any(), gomock.Any()).Return(&uuid.UUID{}, nil)
		mockJwtManager.EXPECT().GenerateTokenPair(gomock.Any()).Return(&jsonwebtoken.TokenPair{}, &jsonwebtoken.AdditionalInfo{}, nil)
		mockSession.EXPECT().AddSession(gomock.Any(), gomock.Any()).Return(nil)

		user := &model.User{
			Email:    "test@example.com",
			Password: "password",
		}

		_, err := s.SignUp(ctx, user)
		assert.NoError(t, err)
	})

	t.Run("error when generating password hash", func(t *testing.T) {
		user := &model.User{
			Email:    "test@example.com",
			Password: string(make([]byte, 1001)), // bcrypt can handle up to 72 bytes
		}

		_, err := s.SignUp(ctx, user)
		assert.Error(t, err)
	})

	t.Run("error when inserting user into the database", func(t *testing.T) {
		user := &model.User{
			Email:    "test@example.com",
			Password: "password",
		}

		mockDatabase.EXPECT().InsertUser(gomock.Any(), gomock.Any()).Return(&uuid.UUID{}, errors.New("database error"))

		_, err := s.SignUp(ctx, user)
		assert.Error(t, err)
	})

	t.Run("error when generating token pair", func(t *testing.T) {
		mockDatabase.EXPECT().InsertUser(gomock.Any(), gomock.Any()).Return(&uuid.UUID{}, nil)
		mockJwtManager.EXPECT().GenerateTokenPair(gomock.Any()).Return(nil, nil, errors.New("token error"))

		user := &model.User{
			Email:    "test@example.com",
			Password: "password",
		}
		_, err := s.SignUp(ctx, user)
		assert.Error(t, err)
	})

	t.Run("error when adding session", func(t *testing.T) {
		mockDatabase.EXPECT().InsertUser(gomock.Any(), gomock.Any()).Return(&uuid.UUID{}, nil)
		mockJwtManager.EXPECT().GenerateTokenPair(gomock.Any()).Return(&jsonwebtoken.TokenPair{}, &jsonwebtoken.AdditionalInfo{}, nil)
		mockSession.EXPECT().AddSession(gomock.Any(), gomock.Any()).Return(errors.New("session error"))

		user := &model.User{
			Email:    "test@example.com",
			Password: "password",
		}

		_, err := s.SignUp(ctx, user)
		assert.Error(t, err)
	})

	t.Run("error when user already exists", func(t *testing.T) {
		mockDatabase.EXPECT().InsertUser(gomock.Any(), gomock.Any()).Return(&uuid.UUID{}, errors.New("user already exists"))

		user := &model.User{
			Email:    "test@example.com",
			Password: "password",
		}

		_, err := s.SignUp(ctx, user)
		assert.Error(t, err)
	})
}
