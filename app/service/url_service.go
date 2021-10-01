package service

import (
	"errors"
	"net/http"
	"os"

	"github.com/moritiza/url-shortener/app/dto"
	"github.com/moritiza/url-shortener/app/entity"
	"github.com/moritiza/url-shortener/app/helper"
	"github.com/moritiza/url-shortener/app/repository"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UrlService interface {
	CreateShortUrl(r *http.Request, url dto.CreateShortUrl) (dto.CreateShortUrl, error)
	Redirect(urlName string) (string, error)
}

// urlService satisfy UrlService interface
type urlService struct {
	logger        logrus.Logger
	urlRepository repository.UrlRepository
}

// NewUrlService creates a new url service with the given dependencies
func NewUrlService(l logrus.Logger, ur repository.UrlRepository) UrlService {
	return &urlService{
		logger:        l,
		urlRepository: ur,
	}
}

// CreateShortUrl do creating short url steps
func (us *urlService) CreateShortUrl(r *http.Request, url dto.CreateShortUrl) (dto.CreateShortUrl, error) {
	// Generate random name for given url
	randString := helper.GetRandomString()

	var ue = entity.Url{
		Title:       url.Title,
		OriginalUrl: url.OriginalUrl,
		UrlName:     randString,
	}

	// Insert new url to urls table
	db := us.urlRepository.Create(ue)
	if db.Error != nil {
		us.logger.Error("Error: ", db.Error)
		return dto.CreateShortUrl{}, db.Error
	}

	// Make short url from generated url name
	url.ShortUrl = os.Getenv("HOST") + ":" + os.Getenv("PORT") + "/" + randString
	return url, nil
}

// Redirect find original url by url unique name and redirect into
func (us *urlService) Redirect(urlName string) (string, error) {
	// Get original url by url name
	url, db := us.urlRepository.GetByName(urlName)
	if db.Error != nil {
		// Check database error type and handle
		if errors.Is(db.Error, gorm.ErrRecordNotFound) {
			return "", errors.New("404")
		}

		return "", db.Error
	}

	return url.OriginalUrl, nil
}
