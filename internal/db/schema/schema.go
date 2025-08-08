package schema

import (
	"gorm.io/gorm"
)

type (
	Run struct {
		gorm.Model
		Job     string `gorm:"not null"`
		Retries uint   `gorm:"default:0;not null"`
	}
)
