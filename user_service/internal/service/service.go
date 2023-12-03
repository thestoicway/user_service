package service

import (
	"context"

	"github.com/thestoicway/backend/user_service/internal/database"
	"go.uber.org/zap"
)

type UserService interface {
	SignIn(ctx context.Context, email string, password string) (string, error)
}

type userService struct {
	logger   *zap.SugaredLogger
	database database.UserDatabase
}

func NewUserService(logger *zap.SugaredLogger, database database.UserDatabase) UserService {
	return &userService{
		logger:   logger,
		database: database,
	}
}
