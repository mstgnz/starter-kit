package handler

import (
	"context"
	"net/http"

	"github.com/mstgnz/starter-kit/api/internal/response"
)

type homeHandler struct{}

func NewHomeHandler() *homeHandler {
	return &homeHandler{}
}

func (h *homeHandler) Home(ctx context.Context, req *any) response.Response {
	return response.Response{
		Code:    http.StatusOK,
		Success: true,
		Message: "home page",
	}
}
