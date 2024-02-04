package service

import (
	"github.com/thestoicway/user_service/internal/profile/database"
	"github.com/thestoicway/user_service/internal/profile/model"
	"go.uber.org/zap"
)

type ProfileService interface {
	UpdateProfile(profile *model.Profile) error
}

type profileServiceImpl struct {
	*ProfileServiceParams
}

// UpdateProfile implements ProfileService.
func (p *profileServiceImpl) UpdateProfile(profile *model.Profile) error {
	profileDB := &model.ProfileDB{
		UserID:    profile.UserID,
		Name:      profile.Name,
		BirthDate: profile.BirthDate,
	}

	if err := p.Database.CreateOrUpdate(profileDB); err != nil {
		return err
	}

	return nil
}

type ProfileServiceParams struct {
	Logger   *zap.SugaredLogger
	Database database.ProfileDatabase
}

func NewProfileService(p *ProfileServiceParams) ProfileService {
	return &profileServiceImpl{
		ProfileServiceParams: p,
	}
}
