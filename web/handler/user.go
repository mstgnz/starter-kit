package handler

import (
	"net/http"
)

type userHandler struct {
}

func NewUserHandler() *userHandler {
	return &userHandler{}
}

func (h *userHandler) Login(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (h *userHandler) Register(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (h *userHandler) Verify(w http.ResponseWriter, r *http.Request) error {
	return nil
}
