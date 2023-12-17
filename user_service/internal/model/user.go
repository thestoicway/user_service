package model

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	customerrors "github.com/thestoicway/backend/custom_errors/custom_errors"
)

type User struct {
	Email    string `json:"email"    validate:"required"`
	Name     string `json:"name"`
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

			return customerrors.NewWrongInputError(fmt.Sprintf("wrong input fields: %v", fields))
		}

		return customerrors.NewWrongInputError(fmt.Sprintf("can't validate user: %v", err.Error()))
	}

	return nil
}

type UserDB struct {
	ID           int
	Email        string
	PasswordHash string
}
