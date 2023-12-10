package service

import (
	"context"

	"github.com/thestoicway/backend/user_service/internal/config"
	"github.com/thestoicway/backend/user_service/internal/database"
	"github.com/thestoicway/backend/user_service/internal/model"
	"go.uber.org/zap"
)

type UserService interface {
	SignIn(ctx context.Context, user *model.User) (string, error)
	SignUp(ctx context.Context, user *model.User) error
}

type userService struct {
	logger   *zap.SugaredLogger
	config   *config.Config
	database database.UserDatabase
}

func NewUserService(logger *zap.SugaredLogger, database database.UserDatabase, config *config.Config) UserService {
	return &userService{
		logger:   logger,
		database: database,
		config:   config,
	}
}
