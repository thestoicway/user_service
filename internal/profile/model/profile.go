package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Profile struct {
	UserID    uuid.UUID
	Name      string    `json:"name"`
	BirthDate time.Time `json:"birth_date"`
}

type ProfileDB struct {
	UserID    uuid.UUID `gorm:"type:uuid;"`
	Name      string
	BirthDate time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (ProfileDB) TableName() string {
	return "profiles"
}
