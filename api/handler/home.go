package handler

import (
	"net/http"

	"github.com/mstgnz/starter-kit/api/internal/load"
	"github.com/mstgnz/starter-kit/api/view/page"
)

type homeHandler struct{}

func NewHomeHandler() *homeHandler {
	return &homeHandler{}
}

func (h *homeHandler) Home(w http.ResponseWriter, r *http.Request) error {
	return load.Render(w, r, page.Home())
}
