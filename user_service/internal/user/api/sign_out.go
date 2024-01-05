package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (h *userHandlerImpl) SignOut(w http.ResponseWriter, r *http.Request, ps httprouter.Params) error {
	refreshToken, err := getRefreshTokenFromQueryParams(r)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	if err = h.service.SignOut(r.Context(), refreshToken); err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	return nil
}
