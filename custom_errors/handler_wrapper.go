package customerrors

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func handleErrorResponse(w http.ResponseWriter, customError *CustomError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(customError.StatusCode())
	resp := NewErrorResponse(customError)
	jsonEncoder := json.NewEncoder(w)
	jsonEncoder.Encode(resp)
}

// HandlerWrapper is a middleware that wraps the handler with a function that
// handles errors and returns a unified response.
func HandlerWrapper(handler func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) error) func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		// Recover from panic and return an internal server error.
		defer func() {
			if err := recover(); err != nil {
				resp := NewInternalServerError(err.(error))
				handleErrorResponse(w, resp)
			}
		}()

		err := handler(w, r, ps)

		if err != nil {
			if err, ok := err.(*CustomError); ok {
				handleErrorResponse(w, err)
				return
			}

			resp := NewInternalServerError(err)

			handleErrorResponse(w, resp)
			return
		}
	}
}
