package service

import (
	"os"

	"github.com/moritiza/url-shortener/app/dto"
	"github.com/moritiza/url-shortener/app/entity"
	"github.com/moritiza/url-shortener/app/helper"
	"github.com/moritiza/url-shortener/app/repository"
	"github.com/sirupsen/logrus"
)

type UrlService interface {
	CreateShortUrl(url dto.CreateShortUrl) (dto.CreateShortUrl, error)
	Redirect(urlName string) (string, error)
	GetShortUrlDetail(urlName string) (dto.GetShortUrlDetail, error)
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
		OriginalUrl: url.OriginalUrl,
	}

	// Insert new url to urls table
	id, db := us.urlRepository.Create(ue)
	if db.Error != nil {
		us.logger.Error("Error: ", db.Error)
		return dto.CreateShortUrl{}, db.Error
	}

	// Create base62 url unique name with inserted url id
	urlName := helper.ToBase62(id)

	// Make short url from generated url name
	url.ShortUrl = os.Getenv("HOST") + ":" + os.Getenv("PORT") + "/" + urlName
	return url, nil
}

// Redirect find original url by url unique name and redirect into
func (us *urlService) Redirect(urlName string) (string, error) {
	// Get url by url name
	url, err := us.getUrl(urlName)
	if err != nil {
		return "", err
	}

	// Increment url click one unit
	db := us.urlRepository.IncrementUrlClick(url.ID)
	if db.Error != nil {
		us.logger.Error("Error: ", db.Error)
		return "", db.Error
	}

	return url.OriginalUrl, nil
}

// GetShortUrlDetail get url by url unique name and return url detail
func (us *urlService) GetShortUrlDetail(urlName string) (dto.GetShortUrlDetail, error) {
	// Get url by url name
	url, err := us.getUrl(urlName)
	if err != nil {
		return dto.GetShortUrlDetail{}, err
	}

	// Prepare and return GetShortUrlDetail DTO
	return dto.GetShortUrlDetail{OriginalUrl: url.OriginalUrl, Click: url.Click}, nil
}

// getUrl get url by url unique name
func (us *urlService) getUrl(urlName string) (entity.Url, error) {
	// Get url id from converting urlName to base10
	id, err := helper.ToBase10(urlName)
	if err != nil {
		us.logger.Error("Error: ", err)
		return entity.Url{}, err
	}

	// Get url by url id
	url, db := us.urlRepository.GetByID(id)
	if db.Error != nil {
		us.logger.Error("Error: ", db.Error)
		return entity.Url{}, db.Error
	}

	return url, nil
}
