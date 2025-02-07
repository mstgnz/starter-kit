package handler

import (
	"net/http"

	"github.com/mstgnz/starter-kit/internal/load"
	"github.com/mstgnz/starter-kit/view/page"
)

type HomeHandler struct{}

func (h *HomeHandler) HomeHandler(w http.ResponseWriter, r *http.Request) error {
	return load.Render(w, r, page.Home())
}
