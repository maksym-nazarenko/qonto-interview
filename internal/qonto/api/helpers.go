package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/maxim-nazarenko/qonto-interview/internal/qonto/core"
)

const (
	HeaderContentType string = "Content-Type"
)

func Respond(w http.ResponseWriter, r *http.Request, content interface{}) {
	RespondCode(w, r, http.StatusOK, content)
}

func RespondCode(w http.ResponseWriter, r *http.Request, code int, content interface{}) {
	w.WriteHeader(code)
	w.Header().Set(HeaderContentType, "application/json")

	b, err := json.Marshal(content)
	if err != nil {
		log.Printf("error marshalling response: %v", err)
		return
	}
	_, err = w.Write(b)
	if err != nil {
		log.Printf("error writing response: %v", err)
		return
	}
}

func handleErrors(w http.ResponseWriter, r *http.Request, err error) {
	wrappedErr := wrapError(err)
	switch {
	case errors.Is(err, core.ErrNotEnoughFunds):
		RespondCode(w, r, http.StatusUnprocessableEntity, wrappedErr)
	case errors.Is(err, core.ErrInvalidCurrency):
		RespondCode(w, r, http.StatusBadRequest, wrappedErr)
	case errors.Is(err, ErrMalformedInput):
		RespondCode(w, r, http.StatusBadRequest, wrappedErr)
	default:
		// you should never expose wild errors in production,
		// ideally, even previous errors must be wrapped in abstract errors
		// without details, unless they are properly handled
		RespondCode(w, r, http.StatusInternalServerError, wrappedErr)
	}
}
