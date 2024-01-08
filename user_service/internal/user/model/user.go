package model

import (
	"fmt"
	"net/mail"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	customerrors "github.com/thestoicway/backend/custom_errors"
	"gorm.io/gorm"
)

type User struct {
	Email    string `json:"email"    validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (u *User) Validate() error {
	validate := validator.New()
	err := validate.Struct(u)

	if err != nil {
		if _, ok := err.(validator.ValidationErrors); ok {
			fields := []string{}

			for _, err := range err.(validator.ValidationErrors) {
				fields = append(fields, err.Field())
			}

			return customerrors.NewWrongInputError(fmt.Errorf("wrong input fields: %v", fields))
		}

		return customerrors.NewWrongInputError(fmt.Errorf("can't validate user: %v", err.Error()))
	}

	// validate if email is not empty and is valid
	_, err = mail.ParseAddress(u.Email)

	if err != nil {
		return customerrors.NewWrongCredentialsError(err)
	}

	return nil
}

type UserDB struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Email        string    `gorm:"unique;not null"`
	PasswordHash string    `gorm:"not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

func (UserDB) TableName() string {
	return "users"
}
