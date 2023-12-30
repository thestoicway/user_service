package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	customerrors "github.com/thestoicway/backend/custom_errors/custom_errors"
)

func (h *userHandler) Refresh(w http.ResponseWriter, r *http.Request, ps httprouter.Params) error {
	// get refresh token from query parameters
	rToken := r.URL.Query().Get("refresh_token")

	if rToken == "" {
		return customerrors.NewWrongInputError(
			"Refresh token should be provided in the request query-parameters",
		)
	}

	newPair, err := h.service.Refresh(r.Context(), rToken)

	if err != nil {
		return err
	}

	// write token pair to response
	jsonEncoder := json.NewEncoder(w)
	jsonEncoder.Encode(customerrors.NewSuccessResponse(newPair))

	return nil
}
