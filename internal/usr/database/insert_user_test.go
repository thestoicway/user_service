package database_test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	customerrors "github.com/thestoicway/custom_errors"
	"github.com/thestoicway/user_service/internal/usr/database"
	"github.com/thestoicway/user_service/internal/usr/model"
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

	t.Run("Expected flow", func(t *testing.T) {
		t.Parallel()

		logger := zaptest.NewLogger(t).Sugar()

		db, mock := newMockDB(t, logger)

		uid := uuid.New()

		mock.ExpectBegin()
		mock.ExpectQuery("INSERT INTO \"users\" (.+) VALUES (.+)").
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uid))
		mock.ExpectCommit()

		email := "qwerty@gmail.com"

		userDB := database.NewUserDatabase(logger, db)

		ctx := context.Background()

		userToInsert := &model.UserDB{
			Email:        email,
			PasswordHash: "password_hash",
		}

		id, err := userDB.InsertUser(ctx, userToInsert)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if id == nil {
			t.Fatalf("Expected id, got nil")
		}

		if *id != uid {
			t.Fatalf("Expected id %v, got %v", uid, id)
		}

		// we make sure that all expectations were met
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("there were unfulfilled expectations: %s", err)
			return
		}
	})

	t.Run("Cancelled context", func(t *testing.T) {
		t.Parallel()

		logger := zaptest.NewLogger(t).Sugar()

		db, _ := newMockDB(t, logger)

		email := "qwerty@gmail.com"

		userDB := database.NewUserDatabase(logger, db)

		ctx, cancel := context.WithCancel(context.Background())

		cancel()

		userToInsert := &model.UserDB{
			Email:        email,
			PasswordHash: "password_hash",
		}

		id, err := userDB.InsertUser(ctx, userToInsert)

		if id != nil {
			t.Fatalf("Expected nil user, got %v", id)
		}

		if err != context.Canceled {
			t.Fatalf("Expected context cancelled error, got %v", err)
		}
	})
}
