package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Adapter struct {
	auth authService
	log  logger
}

func New(auth authService, log logger) *Adapter {
	a := &Adapter{
		auth: auth,
		log:  log,
	}
	return a
}

func (a *Adapter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router := chi.NewRouter()
	router.Use(LoggingMiddleware(a.log))
	router.Route("/api/user", func(r chi.Router) {
		r.Post("/register", a.Register)
		r.Post("/login", a.Login)
	})
	router.ServeHTTP(w, r)
}
