package service

import (
	"errors"
	"os"

	"github.com/moritiza/url-shortener/app/dto"
	"github.com/moritiza/url-shortener/app/entity"
	"github.com/moritiza/url-shortener/app/helper"
	"github.com/moritiza/url-shortener/app/repository"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UrlService interface {
	CreateShortUrl(url dto.CreateShortUrl) (dto.CreateShortUrl, error)
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
func (us *urlService) CreateShortUrl(url dto.CreateShortUrl) (dto.CreateShortUrl, error) {
	// Create url entity from CreateShortUrl DTO
	var ue = entity.Url{
		Title:       url.Title,
		OriginalUrl: url.OriginalUrl,
	}

	// Insert new url to urls table
	id, db := us.urlRepository.Create(ue)
	if db.Error != nil {
		us.logger.Error("Error: ", db.Error)
		return dto.CreateShortUrl{}, db.Error
	}

	// Create base62 urlName with inserted url id
	urlName := helper.ToBase62(id)

	// Make short url from generated url name
	url.ShortUrl = os.Getenv("HOST") + ":" + os.Getenv("PORT") + "/" + urlName
	return url, nil
}

// Redirect find original url by url unique name and redirect into
func (us *urlService) Redirect(urlName string) (string, error) {
	// Get url id from converting urlName to base10
	id, err := helper.ToBase10(urlName)
	if err != nil {
		us.logger.Error("Error: ", err)
		return "", err
	}

	// Get original url by url name
	url, db := us.urlRepository.GetByID(id)
	if db.Error != nil {
		// Check database error type and handle
		if errors.Is(db.Error, gorm.ErrRecordNotFound) {
			return "", errors.New("404")
		}

		return "", db.Error
	}

	return url.OriginalUrl, nil
}
