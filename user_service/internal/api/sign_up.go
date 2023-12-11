package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	customerrors "github.com/thestoicway/backend/custom_errors/custom_errors"
	"github.com/thestoicway/backend/user_service/internal/model"
)

func (h *userHandler) SignUp(w http.ResponseWriter, r *http.Request, ps httprouter.Params) error {
	user := &model.User{}

	err := json.NewDecoder(r.Body).Decode(user)

	if err != nil {
		return customerrors.NewWrongInputError(fmt.Sprintf("can't decode request body: %v", err.Error()))
	}

	if err := user.Validate(); err != nil {
		return err
	}

	if err := h.service.SignUp(r.Context(), user); err != nil {
		return err
	}

	// Make status code 201
	w.WriteHeader(http.StatusCreated)
	return nil
}
