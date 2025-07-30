package schema

import (
	"gorm.io/gorm"
)

type (
	Retry struct {
		gorm.Model
		Job     string `gorm:"unique;not null"`
		Retries uint   `gorm:"default:0;not null"`
	}
)
