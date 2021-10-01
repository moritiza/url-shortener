package repository

import (
	"github.com/moritiza/url-shortener/app/entity"
	"gorm.io/gorm"
)

type UrlRepository interface {
	Create(entity.Url) *gorm.DB
	GetByName(urlName string) (entity.Url, *gorm.DB)
}

// urlRepository satisfy UrlRepository interface
type urlRepository struct {
	db *gorm.DB
}

// NewUrlRepository creates a new url repository with the given dependencies
func NewUrlRepository(db *gorm.DB) UrlRepository {
	return &urlRepository{
		db: db,
	}
}

// Create do insert operation on urls table and return database result
func (ur *urlRepository) Create(url entity.Url) *gorm.DB {
	r := ur.db.Model(entity.Url{}).Create(&url)
	return r
}

// GetByName do read operation on urls table and return founded url with database result
func (ur *urlRepository) GetByName(urlName string) (entity.Url, *gorm.DB) {
	var url entity.Url

	r := ur.db.Model(entity.Url{}).Where("url_name = ?", urlName).First(&url)
	return url, r
}
