package database

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/thestoicway/backend/user_service/internal/user/model"
	customerrors "github.com/thestoicway/custom_errors"
	"gorm.io/gorm"
)

func (db *userDatabaseImpl) InsertUser(ctx context.Context, user *model.UserDB) (userID uuid.UUID, err error) {
	gormDb := db.db

	res := gormDb.Create(user)

	err = res.Error

	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return uuid.UUID{}, customerrors.NewDuplicateEmailError()
		}

		return uuid.UUID{}, err
	}

	return user.ID, nil
}
