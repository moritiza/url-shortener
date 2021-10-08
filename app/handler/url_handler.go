package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/moritiza/url-shortener/app/dto"
	"github.com/moritiza/url-shortener/app/helper"
	"github.com/moritiza/url-shortener/app/service"
	"github.com/moritiza/url-shortener/config"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UrlHandler interface {
	CreateShortUrl(w http.ResponseWriter, r *http.Request)
	Redirect(w http.ResponseWriter, r *http.Request)
	GetShortUrlDetail(w http.ResponseWriter, r *http.Request)
}

// urlHandler is a http.Handler and satisfy UrlHandler interface
type urlHandler struct {
	logger     logrus.Logger
	validator  validator.Validate
	urlService service.UrlService
}

// NewUrlHandler creates a new url handler with the given dependencies
func NewUrlHandler(l logrus.Logger, v validator.Validate, us service.UrlService) UrlHandler {
	return &urlHandler{
		logger:     l,
		validator:  v,
		urlService: us,
	}
}

// CreateShortUrl implements the go http.Handler interface
func (uh *urlHandler) CreateShortUrl(w http.ResponseWriter, r *http.Request) {
	var url dto.CreateShortUrl

	// Decode received data and store them into CreateShortUrl DTO
	err := json.NewDecoder(r.Body).Decode(&url)
	if err != nil {
		// Return 400 error
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate received data
	err = uh.validator.Struct(url)
	if err != nil {
		// Return 400 error with validation errors
		helper.FailureResponse(w, "bad request", config.ValidatorErrors(&uh.validator, err), nil, http.StatusBadRequest)
		return
	}

	// Call url service CreateShortUrl method
	url, err = uh.urlService.CreateShortUrl(url)
	if err != nil {
		// Return 500 error for unhandled errors
		helper.FailureResponse(w, "error", err.Error(), nil, http.StatusInternalServerError)
		return
	}

	// Return Created with header code 201
	helper.SuccessResponse(w, "created", url, true, http.StatusCreated)
}

// Redirect implements the go http.Handler interface
func (uh *urlHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	// Call url service Redirect method
	url, err := uh.urlService.Redirect(mux.Vars(r)["url_name"])
	if err != nil {
		// Check database error type and handle
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Return 404 error
			helper.FailureResponse(w, "not found", "url not found", nil, http.StatusNotFound)
			return
		}

		// Return 500 error for unhandled errors
		helper.FailureResponse(w, "error", err.Error(), nil, http.StatusInternalServerError)
		return
	}

	// Redirect to destination url with header code 303
	http.Redirect(w, r, url, http.StatusSeeOther)
}

// GetShortUrlDetail implements the go http.Handler interface
func (uh *urlHandler) GetShortUrlDetail(w http.ResponseWriter, r *http.Request) {
	// Call url service GetShortUrlDetail method
	detail, err := uh.urlService.GetShortUrlDetail(mux.Vars(r)["url_name"])
	if err != nil {
		// Check database error type and handle
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Return 404 error
			helper.FailureResponse(w, "not found", "url not found", nil, http.StatusNotFound)
			return
		}

		// Return 500 error for unhandled errors
		helper.FailureResponse(w, "error", err.Error(), nil, http.StatusInternalServerError)
		return
	}

	// Return OK with header code 200
	helper.SuccessResponse(w, "ok", detail, true, http.StatusOK)
}
