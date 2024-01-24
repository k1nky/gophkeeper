package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Adapter struct {
	auth   authService
	keeper keeperService
	log    logger
}

func New(auth authService, keeper keeperService, log logger) *Adapter {
	a := &Adapter{
		auth:   auth,
		keeper: keeper,
		log:    log,
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
	router.Route("/api/secret", func(r chi.Router) {
		r.With(AuthorizeMiddleware(a.auth)).Get("/meta/{id}", a.GetSecretMeta)
		r.With(AuthorizeMiddleware(a.auth)).Get("/{id}", a.GetSecretData)
		r.With(AuthorizeMiddleware(a.auth)).Put("/", a.PutSecret)
		r.With(AuthorizeMiddleware(a.auth)).Put("/file", a.PutSecretFile)
	})
	router.ServeHTTP(w, r)
}

func (a *Adapter) writeJSON(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		return err
	}
	return nil
}
