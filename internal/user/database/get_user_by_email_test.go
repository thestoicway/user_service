package database_test

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/thestoicway/backend/user_service/internal/user/database"
	customerrors "github.com/thestoicway/custom_errors"
	"go.uber.org/zap/zaptest"
)

func TestGetUserByEmail(t *testing.T) {
	t.Parallel()

	t.Run("ExistingEmail", func(t *testing.T) {
		t.Parallel()

		logger := zaptest.NewLogger(t).Sugar()

		db, mock := newMockDB(t, logger)

		email := "qwerty@gmail.com"

		rows := sqlmock.NewRows([]string{"id", "email", "password_hash"}).
			AddRow(uuid.New(), email, "password_hash")

		mock.ExpectQuery("SELECT (.+) FROM \"users\" WHERE email = (.+)").
			WithArgs(email).
			WillReturnRows(rows)

		userDB := database.NewUserDatabase(logger, db)

		ctx := context.Background()

		user, err := userDB.GetUserByEmail(ctx, email)

		if err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}

		if user == nil {
			t.Fatalf("User is nil")
		}

		if user.Email != email {
			t.Fatalf("Emails do not match")
		}
	})

	t.Run("NonExistingEmail", func(t *testing.T) {
		t.Parallel()

		logger := zaptest.NewLogger(t).Sugar()

		db, mock := newMockDB(t, logger)

		rows := sqlmock.NewRows([]string{"id", "email", "password_hash"})

		mock.ExpectQuery("SELECT (.+) FROM \"users\" WHERE email = (.+)").
			WillReturnRows(rows)

		userDB := database.NewUserDatabase(logger, db)

		ctx := context.Background()

		user, err := userDB.GetUserByEmail(ctx, "qwerty@gmail.com")

		if user != nil {
			t.Fatalf("User is not nil")
		}

		if err == nil {
			t.Fatalf("Expected error, got nil")
		}

		if err, ok := err.(*customerrors.CustomError); ok {
			if err.Code != customerrors.ErrWrongCredentials {
				t.Fatalf("Expected error code: %d, got: %d", customerrors.ErrWrongCredentials, err.Code)
			}
		} else {
			t.Fatalf("Expected error to be of type *CustomError, got: %T", err)
		}
	})

	t.Run("ComplexEmail", func(t *testing.T) {
		t.Parallel()

		logger := zaptest.NewLogger(t).Sugar()

		db, mock := newMockDB(t, logger)

		email := "very+very.uncommon.part.with+plus.and.dot@subdomain.example.co.uk"

		rows := sqlmock.NewRows([]string{"id", "email", "password_hash"}).
			AddRow(uuid.New(), email, "password_hash")

		mock.ExpectQuery("SELECT (.+) FROM \"users\" WHERE email = (.+)").
			WithArgs(email).
			WillReturnRows(rows)

		userDB := database.NewUserDatabase(logger, db)

		ctx := context.Background()

		user, err := userDB.GetUserByEmail(ctx, email)

		if err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}

		if user == nil {
			t.Fatalf("User is nil")
		}

		if user.Email != email {
			t.Fatalf("Emails do not match")
		}
	})

	t.Run("InvalidEmail", func(t *testing.T) {
		t.Parallel()

		logger := zaptest.NewLogger(t).Sugar()

		db, mock := newMockDB(t, logger)

		rows := sqlmock.NewRows([]string{"id", "email", "password_hash"}).
			AddRow(uuid.New(), "qwerty", "password_hash")

		mock.ExpectQuery("SELECT (.+) FROM \"users\" WHERE email = (.+)").
			WillReturnRows(rows)

		userDB := database.NewUserDatabase(logger, db)

		ctx := context.Background()

		user, err := userDB.GetUserByEmail(ctx, "qwerty")

		if user != nil {
			t.Fatalf("User is not nil")
		}

		if err == nil {
			t.Fatalf("Expected error, got nil")
		}

		if err, ok := err.(*customerrors.CustomError); ok {
			if err.Code != customerrors.ErrWrongCredentials {
				t.Fatalf("Expected error code: %d, got: %d", customerrors.ErrWrongCredentials, err.Code)
			}
		} else {
			t.Fatalf("Expected error to be of type *CustomError, got: %T", err)
		}
	})

	t.Run("Context Cancelled", func(t *testing.T) {
		t.Parallel()

		logger := zaptest.NewLogger(t).Sugar()

		db, mock := newMockDB(t, logger)

		email := "qwerty@gmail.com"

		rows := sqlmock.NewRows([]string{"id", "email", "password_hash"}).
			AddRow(uuid.New(), email, "password_hash")

		mock.ExpectQuery("SELECT (.+) FROM \"users\" WHERE email = (.+)").
			WithArgs(email).
			WillReturnRows(rows)

		userDB := database.NewUserDatabase(logger, db)

		ctx, cancel := context.WithCancel(context.Background())

		cancel()

		user, err := userDB.GetUserByEmail(ctx, email)

		if user != nil {
			t.Fatalf("User is not nil")
		}

		if err == nil {
			t.Fatalf("Expected error, got nil")
		}

		expectedError := context.Canceled

		if !errors.Is(err, expectedError) {
			t.Fatalf("Expected error: %s, got: %s", expectedError, err)
		}
	})
}
