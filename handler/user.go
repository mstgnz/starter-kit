package handler

import (
	"net/http"
)

type UserHandler struct {
}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (h *UserHandler) Verify(w http.ResponseWriter, r *http.Request) error {
	return nil
}
