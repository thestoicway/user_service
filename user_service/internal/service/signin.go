package service

import (
	"context"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

func (svc *userService) SignIn(ctx context.Context, email string, password string) (string, error) {
	_, err := svc.database.GetUserByEmail(email)

	if err != nil {
		return "", fmt.Errorf("can't get user by email: %v", err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{})

	tokenString, err := token.SignedString([]byte("secret"))

	if err != nil {
		return "", fmt.Errorf("can't sign token: %v", err)
	}

	return tokenString, nil
}
