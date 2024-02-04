package database

import (
	"context"

	"github.com/google/uuid"
	"github.com/thestoicway/user_service/internal/usr/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserDatabase interface {
	GetUserByEmail(ctx context.Context, email string) (*model.UserDB, error)
	InsertUser(ctx context.Context, user *model.UserDB) (userID *uuid.UUID, err error)
}

type userDatabaseImpl struct {
	logger *zap.SugaredLogger
	db     *gorm.DB
}

func NewUserDatabase(logger *zap.SugaredLogger, db *gorm.DB) UserDatabase {
	return &userDatabaseImpl{
		logger: logger,
		db:     db,
	}
}
