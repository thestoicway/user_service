package database

import (
	"context"

	"github.com/thestoicway/backend/user_service/internal/model"
	"go.uber.org/zap"
)

type UserDatabase interface {
	GetUserByEmail(ctx context.Context, email string) (*model.UserDB, error)
	InsertUser(ctx context.Context, user *model.UserDB) (userID int, err error)
}

type userDatabase struct {
	logger *zap.SugaredLogger
	db     *DB
}

func NewUserDatabase(logger *zap.SugaredLogger, db *DB) UserDatabase {
	return &userDatabase{
		logger: logger,
		db:     db,
	}
}
