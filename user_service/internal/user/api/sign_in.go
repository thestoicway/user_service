package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	customerrors "github.com/thestoicway/backend/custom_errors"
	"github.com/thestoicway/backend/user_service/internal/user/model"
)

func (h *userHandlerImpl) SignIn(w http.ResponseWriter, r *http.Request, ps httprouter.Params) error {

	user := &model.User{}

	err := json.NewDecoder(r.Body).Decode(user)

	if err != nil {
		return customerrors.NewWrongInputError(
			fmt.Sprintf("can't decode request body: %v", err.Error()),
		)
	}

	if err := user.Validate(); err != nil {
		return err
	}

	token, err := h.service.SignIn(r.Context(), user)

	if err != nil {
		return err
	}

	resp := customerrors.NewSuccessResponse(token)
	return sendJSONResponse(w, resp)
}

func sendJSONResponse(w http.ResponseWriter, resp interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(resp)
}
