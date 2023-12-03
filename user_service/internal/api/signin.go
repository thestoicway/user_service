package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (h *userHandler) SignIn(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	jwt, err := h.service.SignIn(r.Context(), "email", "password")

	if err != nil {
		h.logger.Errorf("can't sign in: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("can't sign in"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(jwt))
}
