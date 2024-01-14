package service

import (
	"context"
	"fmt"

	customerrors "github.com/thestoicway/custom_errors"
	"github.com/thestoicway/user_service/internal/jsonwebtoken"
	"github.com/thestoicway/user_service/internal/model"
	"golang.org/x/crypto/bcrypt"
)

func (s *userServiceImpl) SignIn(ctx context.Context, user *model.User) (tokenPair *jsonwebtoken.TokenPair, err error) {
	userDb, err := s.Database.GetUserByEmail(ctx, user.Email)

	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(userDb.PasswordHash), []byte(user.Password))

	if err != nil {
		return nil, customerrors.NewWrongCredentialsError(err)
	}

	pair, info, err := s.JwtManager.GenerateTokenPair(userDb.ID)

	if err != nil {
		return nil, fmt.Errorf("can't sign token: %v", err)
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
