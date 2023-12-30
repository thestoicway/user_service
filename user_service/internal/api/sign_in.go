package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	customerrors "github.com/thestoicway/backend/custom_errors/custom_errors"
	"github.com/thestoicway/backend/user_service/internal/model"
)

func (h *userHandler) SignIn(w http.ResponseWriter, r *http.Request, ps httprouter.Params) error {

	user := &model.User{}

	err := json.NewDecoder(r.Body).Decode(user)

	if err != nil {
		return customerrors.NewWrongInputError(fmt.Sprintf("can't decode request body: %v", err.Error()))
	}

	if err := user.Validate(); err != nil {
		return err
	}

	jwt, err := h.service.SignIn(r.Context(), user)

	if err != nil {
		return err
	}

	resp := customerrors.NewSuccessResponse(jwt)

	jsonEncoder := json.NewEncoder(w)
	jsonEncoder.Encode(resp)

	return nil
}
