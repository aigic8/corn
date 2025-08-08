package db

import (
	"fmt"

	"github.com/aigic8/corn/internal/db/models"
	"github.com/aigic8/corn/internal/db/schema"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DbClient = gorm.DB

var ErrNotExist = gorm.ErrRecordNotFound

type (
	Db struct {
		DbAddr string
		c      *DbClient
		Retry  *models.RunModel
	}
)

func NewDb(dbAddr string) (*Db, error) {
	c, err := gorm.Open(sqlite.Open(dbAddr), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("connecting to db: %w", err)
	}

	retryModel := models.NewRunModel(c)

	return &Db{DbAddr: dbAddr, c: c, Retry: retryModel}, nil
}

func (db *Db) Init() error {
	return db.c.AutoMigrate(&schema.Run{})
}
