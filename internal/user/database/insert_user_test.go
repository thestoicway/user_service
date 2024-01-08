package database_test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/thestoicway/backend/user_service/internal/user/database"
	"github.com/thestoicway/backend/user_service/internal/user/model"
	customerrors "github.com/thestoicway/custom_errors"
	"go.uber.org/zap/zaptest"
	"gorm.io/gorm"
)

func TestInserUser(t *testing.T) {
	t.Parallel()

	t.Run("ExistingEmail", func(t *testing.T) {
		t.Parallel()

		logger := zaptest.NewLogger(t).Sugar()

		db, mock := newMockDB(t, logger)

		mock.ExpectBegin()
		mock.ExpectQuery("INSERT INTO \"users\" (.+) VALUES (.+)").
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnError(gorm.ErrDuplicatedKey)
		mock.ExpectRollback()

		email := "qwerty@gmail.com"

		userDB := database.NewUserDatabase(logger, db)

		ctx := context.Background()

		userToInsert := &model.UserDB{
			Email:        email,
			PasswordHash: "password_hash",
		}

		id, err := userDB.InsertUser(ctx, userToInsert)

		if id != nil {
			t.Fatalf("Expected nil user, got %v", id)
		}

		// we make sure that all expectations were met
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("there were unfulfilled expectations: %s", err)
			return
		}

		if customError, ok := err.(*customerrors.CustomError); ok {
			if customError.Code != customerrors.ErrDuplicateEmail {
				t.Fatalf("Expected duplicate email error, got %v", customError)
			}
		} else {
			t.Fatalf("Expected custom error, got %v", err)
		}
	})
}
