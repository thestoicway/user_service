package database

import (
	"github.com/thestoicway/backend/user_service/internal/model"
	"go.uber.org/zap"
)

type UserDatabase interface {
	GetUserByEmail(email string) (model.User, error)
}

type userDatabase struct {
	logger *zap.SugaredLogger
}

func NewUserDatabase(logger *zap.SugaredLogger) UserDatabase {
	return &userDatabase{
		logger: logger,
	}
}
