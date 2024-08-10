package handler

import (
	"net/http"
)

type UserHandler struct {
}

func (h *UserHandler) LoginHandler(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (h *UserHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) error {
	return nil
}
