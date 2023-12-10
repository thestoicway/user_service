package database

import "github.com/thestoicway/backend/user_service/internal/model"

func (db *userDatabase) InsertUser(user *model.UserDB) error {
	db.logger.Infof("Inserting user: %v", user)

	// TODO: implement this method
	return nil
}
