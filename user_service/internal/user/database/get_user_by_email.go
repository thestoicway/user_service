package database

import (
	"context"
	"errors"

	"net/mail"

	customerrors "github.com/thestoicway/backend/custom_errors"
	"github.com/thestoicway/backend/user_service/internal/user/model"
	"gorm.io/gorm"
)

func (db *userDatabaseImpl) GetUserByEmail(context context.Context, email string) (*model.UserDB, error) {
	// validate if email is not empty and is valid
	_, err := mail.ParseAddress(email)

	if err != nil {
		return nil, customerrors.NewWrongCredentialsError()
	}

	gormDb := db.db

	var user model.UserDB

	err = gormDb.Where("email = ?", email).First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, customerrors.NewWrongCredentialsError()
		}

		return nil, err
	}

	return &user, nil
}
