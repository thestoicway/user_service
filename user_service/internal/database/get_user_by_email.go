package database

import (
	"github.com/thestoicway/backend/user_service/internal/model"
)

func (db *userDatabase) GetUserByEmail(email string) (*model.UserDB, error) {

	return &model.UserDB{
		Email:        "qwerty@gmail.com",
		Name:         "qwerty",
		PasswordHash: "$2a$10$lr/ilYBVAYQHQv/Z310uwujBKLXp3eoB3ZujlSODkIggvQ3xbKyc6",
	}, nil
}
