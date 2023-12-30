package service

import (
	"context"
	"fmt"

	"github.com/thestoicway/backend/user_service/internal/jsonwebtoken"
	"github.com/thestoicway/backend/user_service/internal/user/model"
	"golang.org/x/crypto/bcrypt"
)

func (s *userService) SignUp(ctx context.Context, user *model.User) (tokenPair *jsonwebtoken.TokenPair, err error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		return nil, fmt.Errorf("can't generate password hash: %v", err)
	}

	userDB := &model.UserDB{
		Email:        user.Email,
		PasswordHash: string(passwordHash),
	}

	id, err := s.Database.InsertUser(ctx, userDB)

	if err != nil {
		return nil, err
	}

	pair, info, err := s.JwtManager.GenerateTokenPair(id)

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
