package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	customerrors "github.com/thestoicway/backend/custom_errors"
)

func (h *userHandlerImpl) Refresh(w http.ResponseWriter, r *http.Request, ps httprouter.Params) error {
	refreshToken, err := getRefreshTokenFromQueryParams(r)
	if err != nil {
		return err
	}

	newPair, err := h.service.Refresh(r.Context(), refreshToken)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	json.NewEncoder(w).Encode(customerrors.NewSuccessResponse(newPair))

	return nil
}

func getRefreshTokenFromQueryParams(r *http.Request) (string, error) {
	refreshToken := r.URL.Query().Get("refresh_token")
	if refreshToken == "" {
		return "", customerrors.NewWrongInputError(
			"Refresh token should be provided in the request query-parameters",
		)
	}
	return refreshToken, nil
}
