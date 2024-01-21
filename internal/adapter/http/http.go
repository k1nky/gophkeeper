package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Adapter struct {
	log logger
}

// /api/user/register,login
// GET /api/vault/{id}
// GET /api/vault
// POST /api/vault

func (a *Adapter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router := chi.NewRouter()
	router.Use(LoggingMiddleware(a.log))
	router.Route("/api/user", func(r chi.Router) {
		r.Post("/register", a.Register)
		// r.Post("/login", a.Login)
		// r.With(AuthorizeMiddleware(a.auth)).Get("/balance", a.GetBalance)
		// r.With(AuthorizeMiddleware(a.auth)).Post("/orders", a.NewOrder)
	})
	router.ServeHTTP(w, r)
}
