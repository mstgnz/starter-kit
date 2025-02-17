package config

import (
	"log"
	"net/http"
)

type HttpHandler func(w http.ResponseWriter, r *http.Request) error

func Catch(h HttpHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			log.Println("HTTP Handler Error", "err", err.Error(), "path", r.URL.Path)
		}
	}
}
