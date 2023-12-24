package database

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/thestoicway/backend/user_service/internal/config"
	"github.com/thestoicway/backend/user_service/internal/model"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type UserDatabase interface {
	GetUserByEmail(ctx context.Context, email string) (*model.UserDB, error)
	InsertUser(ctx context.Context, user *model.UserDB) (userID uuid.UUID, err error)
}

type userDatabase struct {
	logger *zap.SugaredLogger
	db     *gorm.DB
}

func NewUserDatabase(logger *zap.SugaredLogger, db *gorm.DB) UserDatabase {
	return &userDatabase{
		logger: logger,
		db:     db,
	}
}

func OpenDB(config *config.Config) (*gorm.DB, error) {
	// "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable"
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		config.PostgresDatabase.PostgresHost,
		config.PostgresDatabase.PostgresUser,
		config.PostgresDatabase.PostgresPass,
		config.PostgresDatabase.PostgresDB,
		config.PostgresDatabase.PostgresPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		TranslateError: true,
	})

	if err != nil {
		return nil, err
	}

	if err := migrate(db); err != nil {
		return nil, err
	}

	return db, err
}

func migrate(db *gorm.DB) error {

	// Create "users" table
	err := db.Migrator().CreateTable(&model.UserDB{})

	if err != nil {
		return err
	}

	return nil
}
