package handlers

import (
	"errors"
	"net/http"

	"github.com/go-chi/render"
)

func respondError(w http.ResponseWriter, r *http.Request, err error) {
	if errors.Is(err, ErrInternalServer) {
		w.WriteHeader(http.StatusInternalServerError)
	} else if errors.Is(err, ErrInvalidCredentials) {
		w.WriteHeader(http.StatusUnauthorized)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}

	render.JSON(w, r, errorResponse{
		Response: Response{
			Status: statusError,
		},
		Error: err.Error(),
	})
}

func respondOk(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, Response{
		Status: statusSuccess,
	})
}

type Response struct {
	Status string `json:"status"`
}

type errorResponse struct {
	Response
	Error string `json:"error"`
}
