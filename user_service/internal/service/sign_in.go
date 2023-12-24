package service

import (
	"context"
	"fmt"

	customerrors "github.com/thestoicway/backend/custom_errors/custom_errors"
	"github.com/thestoicway/backend/user_service/internal/jsonwebtoken"
	"github.com/thestoicway/backend/user_service/internal/model"
	"golang.org/x/crypto/bcrypt"
)

func (svc *userService) SignIn(ctx context.Context, user *model.User) (tokenPair *jsonwebtoken.TokenPair, err error) {
	userDb, err := svc.Database.GetUserByEmail(ctx, user.Email)

	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(userDb.PasswordHash), []byte(user.Password))

	if err != nil {
		return nil, customerrors.NewWrongCredentialsError()
	}

	pair, err := svc.JwtManager.GenerateTokenPair(userDb.ID)

	if err != nil {
		return nil, fmt.Errorf("can't sign token: %v", err)
	}

	return pair, nil
}
