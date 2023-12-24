package service

import (
	"context"
	"fmt"

	jsonwebtoken "github.com/thestoicway/backend/user_service/internal/jwt"
	"github.com/thestoicway/backend/user_service/internal/model"
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

	pair, err := s.JwtManager.GenerateTokenPair(id)

	if err != nil {
		return nil, err
	}

	return pair, nil
}
