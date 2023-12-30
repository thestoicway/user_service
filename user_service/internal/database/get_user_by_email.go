package database

import (
	"context"
	"errors"

	customerrors "github.com/thestoicway/backend/custom_errors/custom_errors"
	"github.com/thestoicway/backend/user_service/internal/model"
	"gorm.io/gorm"
)

func (db *userDatabase) GetUserByEmail(context context.Context, email string) (*model.UserDB, error) {
	gormDb := db.db

	var user model.UserDB

	err := gormDb.Where("email = ?", email).First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, customerrors.NewWrongCredentialsError()
		}

		return nil, err
	}

	return &user, nil
}
