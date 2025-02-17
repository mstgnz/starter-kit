package handler

import (
	"net/http"

	"github.com/mstgnz/starter-kit/web/infra/render"
	"github.com/mstgnz/starter-kit/web/view/page"
)

type homeHandler struct{}

func NewHomeHandler() *homeHandler {
	return &homeHandler{}
}

func (h *homeHandler) Home(w http.ResponseWriter, r *http.Request) error {
	return render.Render(w, r, page.Home())
}
