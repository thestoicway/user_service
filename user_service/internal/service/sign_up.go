package service

import (
	"context"
	"fmt"

	"github.com/thestoicway/backend/user_service/internal/model"
	"golang.org/x/crypto/bcrypt"
)

func (s *userService) SignUp(ctx context.Context, user *model.User) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		return fmt.Errorf("can't generate password hash: %v", err)
	}

	userDB := &model.UserDB{
		Email:        user.Email,
		PasswordHash: string(passwordHash),
	}

	err = s.database.InsertUser(userDB)

	if err != nil {
		return err
	}

	return nil
}
