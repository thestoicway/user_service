package database

import (
	"github.com/thestoicway/backend/user_service/internal/model"
)

func (db *userDatabase) GetUserByEmail(email string) (model.User, error) {
	return model.User{}, nil
}
