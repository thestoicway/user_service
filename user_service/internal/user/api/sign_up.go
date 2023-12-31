package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	customerrors "github.com/thestoicway/backend/custom_errors"
	"github.com/thestoicway/backend/user_service/internal/user/model"
)

func (h *userHandlerImpl) SignUp(w http.ResponseWriter, r *http.Request, ps httprouter.Params) error {
	user := &model.User{}

	err := json.NewDecoder(r.Body).Decode(user)

	if err != nil {
		return customerrors.NewWrongInputError(fmt.Sprintf("can't decode request body: %v", err.Error()))
	}

	if err := user.Validate(); err != nil {
		return err
	}

	pair, err := h.service.SignUp(r.Context(), user)

	if err != nil {
		return err
	}

	resp := customerrors.NewSuccessResponse(pair)

	jsonEncoder := json.NewEncoder(w)

	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	// Make status code 201
	w.WriteHeader(http.StatusCreated)

	// Write response body
	jsonEncoder.Encode(resp)
	return nil
}
