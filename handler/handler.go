package handler

import (
	"context"
	"net/http"

	"github.com/prodhe/matlistan/template"
)

type handler struct {
	mux *http.ServeMux
}

func New() *handler {
	h := &handler{
		mux: http.NewServeMux(),
	}

	h.mux.Handle("/static/", h.noCache(http.StripPrefix("/static", http.FileServer(http.Dir("./assets/")))))

	h.mux.HandleFunc("/help", h.help)

	h.mux.HandleFunc("/", h.sessionValidate(h.index))

	return h
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

func (h *handler) index(w http.ResponseWriter, r *http.Request) {
	data := template.Fields{
		"Authenticated": true,
	}
	template.Render(w, "index", data)
}

func (h *handler) help(w http.ResponseWriter, r *http.Request) {
	data := template.Fields{
		"Authenticated": true,
	}
	template.Render(w, "help", data)
}

func (h *handler) sessionValidate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = context.WithValue(ctx, "pid", 1)

		next(w, r.WithContext(ctx))
	}
}

func (h *handler) noCache(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate") // HTTP 1.1.
		w.Header().Set("Pragma", "no-cache")                                   // HTTP 1.0.
		w.Header().Set("Expires", "0")                                         // Proxies.
		next.ServeHTTP(w, r)
	}
}
