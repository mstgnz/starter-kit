package handler

import (
	"context"
	"net/http"

	"github.com/mstgnz/starter-kit/api/infra/response"
	"github.com/mstgnz/starter-kit/api/model"
)

type userHandler struct {
}

func NewUserHandler() *userHandler {
	return &userHandler{}
}

func (h *userHandler) Login(ctx context.Context, req *model.Login) response.Response {
	return response.Response{
		Code:    http.StatusOK,
		Success: true,
		Message: "Login successful",
	}
}

func (h *userHandler) Register(ctx context.Context, req *model.Register) response.Response {
	return response.Response{
		Code:    http.StatusOK,
		Success: true,
		Message: "Register successful",
	}
}

func (h *userHandler) Verify(ctx context.Context, _ *any) response.Response {
	return response.Response{
		Code:    http.StatusOK,
		Success: true,
		Message: "Verify successful",
	}
}
