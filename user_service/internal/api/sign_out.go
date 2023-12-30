package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (h *userHandler) SignOut(w http.ResponseWriter, r *http.Request, ps httprouter.Params) error {
	return nil
}
