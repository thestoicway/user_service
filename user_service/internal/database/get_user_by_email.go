package database

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/thestoicway/backend/user_service/internal/model"
)

func (db *userDatabase) GetUserByEmail(context context.Context, email string) (*model.UserDB, error) {
	pool := db.db.DbPool

	user := &model.UserDB{}

	query, args, err := sq.Select("user_id", "email", "password").
		From("users").
		Where(sq.Eq{"email": email}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, err
	}

	row := pool.QueryRow(context, query, args...)

	err = row.Scan(&user.ID, &user.Email, &user.PasswordHash)

	if err != nil {
		return nil, err
	}

	return user, nil
}
