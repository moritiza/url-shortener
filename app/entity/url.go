package entity

import (
	"gorm.io/gorm"
)

// Url entity known as Url model and create urls table
type Url struct {
	ID          uint64 `gorm:"column:id;primaryKey;autoIncrement"`
	Title       string `gorm:"<-;column:title;type:text;not null;comment:title of the URL;example:google"`
	OriginalUrl string `gorm:"<-;column:original_url;type:text;not null;comment:the Url that we want to redirect to;example:http://google.com"`
	UrlName     string `gorm:"<-;column:url_name;type:text;not null;unique;comment:name of the Url that we want to redirect to;example:xwyut2E"`
	gorm.Model
}
