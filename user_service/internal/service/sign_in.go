package service

import (
	"context"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	customerrors "github.com/thestoicway/backend/custom_errors/custom_errors"
	"github.com/thestoicway/backend/user_service/internal/model"
	"golang.org/x/crypto/bcrypt"
)

func (svc *userService) SignIn(ctx context.Context, user *model.User) (string, error) {
	userDb, err := svc.database.GetUserByEmail(user.Email)

	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(userDb.PasswordHash), []byte(user.Password))

	if err != nil {
		return "", customerrors.NewWrongCredentialsError()
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{})

	tokenString, err := token.SignedString([]byte(svc.config.JwtSecret))

	if err != nil {
		return "", fmt.Errorf("can't sign token: %v", err)
	}

	return tokenString, nil
}
