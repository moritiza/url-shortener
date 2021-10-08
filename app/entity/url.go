package entity

import (
	"gorm.io/gorm"
)

// Url entity known as Url model and create urls table
type Url struct {
	ID          uint64 `gorm:"column:id;primaryKey;autoIncrement"`
	OriginalUrl string `gorm:"<-;column:original_url;type:varchar(255);not null;comment:the Url that we want to redirect to;example:http://google.com"`
	gorm.Model
}
