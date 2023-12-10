package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/thestoicway/backend/user_service/internal/model"
)

func (h *userHandler) SignUp(w http.ResponseWriter, r *http.Request, ps httprouter.Params) error {
	user := &model.User{}

	err := json.NewDecoder(r.Body).Decode(user)

	if err != nil {
		return err
	}

	err = h.service.SignUp(r.Context(), user)

	if err != nil {
		return err
	}

	// Make status code 201
	w.WriteHeader(http.StatusCreated)
	return nil
}
