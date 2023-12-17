package database

import (
	"context"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgconn"
	customerrors "github.com/thestoicway/backend/custom_errors/custom_errors"
	"github.com/thestoicway/backend/user_service/internal/model"
)

func (db *userDatabase) InsertUser(ctx context.Context, user *model.UserDB) (userID int, err error) {
	pool := db.db.DbPool

	query, args, err := sq.Insert("users").
		Columns("email", "password").
		Values(user.Email, user.PasswordHash).
		Suffix("RETURNING \"user_id\"").
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return 0, err
	}

	row := pool.QueryRow(ctx, query, args...)

	err = row.Scan(&userID)

	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {

			// if unique constraint violation
			if pgErr.Code == "23505" {
				return 0, customerrors.NewDuplicateEmailError()
			}
		}

		return 0, err
	}

	return userID, nil
}
