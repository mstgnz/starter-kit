package handler

import (
	"context"
	"net/http"

	"github.com/mstgnz/starter-kit/api/infra/response"
)

type UserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type userHandler struct {
}

func NewUserHandler() *userHandler {
	return &userHandler{}
}

func (h *userHandler) Login(ctx context.Context, req *UserRequest) response.Response {
	return response.Response{
		Code:    http.StatusOK,
		Success: true,
		Message: "Login successful",
	}
}

func (h *userHandler) Register(ctx context.Context, req *UserRequest) response.Response {
	return response.Response{
		Code:    http.StatusOK,
		Success: true,
		Message: "Register successful",
	}
}

func (h *userHandler) Verify(ctx context.Context, req *UserRequest) response.Response {
	return response.Response{
		Code:    http.StatusOK,
		Success: true,
		Message: "Verify successful",
	}
}
