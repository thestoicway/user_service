package database

import (
	"github.com/thestoicway/user_service/internal/profile/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ProfileDatabase interface {
	// CreateOrUpdate creates or updates a profile
	// So when there is no profile for the user, it will create a new profile
	// When there is a profile for the user, it will update the existing profile
	CreateOrUpdate(profile *model.ProfileDB) error
}

type profileDatabaseImpl struct {
	logger *zap.SugaredLogger
	db     *gorm.DB
}

// CreateOrUpdate implements ProfileDatabase.
func (p *profileDatabaseImpl) CreateOrUpdate(profile *model.ProfileDB) error {
	if err := p.db.FirstOrCreate(profile, profile).Error; err != nil {
		return err
	}

	return nil
}

func NewProfileDatabase(logger *zap.SugaredLogger, db *gorm.DB) ProfileDatabase {
	return &profileDatabaseImpl{
		logger: logger,
		db:     db,
	}
}
