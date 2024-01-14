package database

import (
	"context"
	"errors"

	"github.com/google/uuid"
	customerrors "github.com/thestoicway/custom_errors"
	"github.com/thestoicway/user_service/internal/model"
	"gorm.io/gorm"
)

func (db *userDatabaseImpl) InsertUser(ctx context.Context, user *model.UserDB) (userID *uuid.UUID, err error) {
	gormDb := db.db.WithContext(ctx)

	res := gormDb.Create(user)

	err = res.Error

	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, customerrors.NewDuplicateEmailError()
		}

		return nil, err
	}

	return &user.ID, nil
}
