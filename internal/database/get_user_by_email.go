package database

import (
	"context"
	"errors"

	customerrors "github.com/thestoicway/custom_errors"

	"github.com/thestoicway/user_service/internal/model"
	"gorm.io/gorm"
)

func (db *userDatabaseImpl) GetUserByEmail(context context.Context, email string) (*model.UserDB, error) {
	gormDb := db.db.WithContext(context)

	var user model.UserDB

	err := gormDb.Where("email = ?", email).First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, customerrors.NewWrongCredentialsError(err)
		}

		return nil, err
	}

	return &user, nil
}
