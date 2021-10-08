package repository

import (
	"fmt"

	"github.com/moritiza/url-shortener/app/entity"
	"gorm.io/gorm"
)

type UrlRepository interface {
	Create(entity.Url) (uint64, *gorm.DB)
	GetByID(id uint64) (entity.Url, *gorm.DB)
	IncrementUrlClick(id uint64) *gorm.DB
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
func (ur *urlRepository) Create(url entity.Url) (uint64, *gorm.DB) {
	r := ur.db.Model(entity.Url{}).Create(&url)
	return url.ID, r
}

// GetByID do read operation on urls table, find url by id and return founded url with database result
func (ur *urlRepository) GetByID(id uint64) (entity.Url, *gorm.DB) {
	var url entity.Url

	r := ur.db.Model(entity.Url{}).Where("id = ?", id).First(&url)
	return url, r
}

// IncrementUrlClick increment url click one unit
func (ur *urlRepository) IncrementUrlClick(id uint64) *gorm.DB {
	r := ur.db.Exec(
		"UPDATE \"urls\" SET click=click+1 WHERE id = " +
			fmt.Sprint(id) + " AND \"urls\".\"deleted_at\" IS NULL",
	)
	return r
}
