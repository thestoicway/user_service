package customerrors

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// HandlerWrapper is a middleware that wraps the handler with a function that
// handles errors and returns a unified response.
func HandlerWrapper(handler func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) error) func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		err := handler(w, r, ps)

		if err != nil {
			if err, ok := err.(*CustomError); ok {
				jsonEncoder := json.NewEncoder(w)
				w.WriteHeader(err.StatusCode)
				resp := NewErrorResponse(err)
				jsonEncoder.Encode(resp)
				return
			}

			resp := NewInternalServerException(err)

			w.WriteHeader(http.StatusInternalServerError)
			jsonEncoder := json.NewEncoder(w)
			jsonEncoder.Encode(resp)
			return
		}
	}
}
