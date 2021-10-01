package core

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/moritiza/url-shortener/app/helper"
	"github.com/moritiza/url-shortener/config"
)

// router create gorilla mux router and define routes
func Router(cfg config.Config) *mux.Router {
	r := mux.NewRouter()
	d := PrepareDependensies(cfg)

	r.HandleFunc("/{url_name}", d.Handlers.UrlHandler.Redirect).Methods(http.MethodGet)

	// Create group routes
	s := r.PathPrefix("/api").Subrouter()

	s.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		helper.SuccessResponse(w, "ok", "pong", true, http.StatusOK)
	}).Methods(http.MethodGet)

	s.HandleFunc("/create-url", d.Handlers.UrlHandler.CreateShortUrl).Methods(http.MethodPost)

	return r
}
