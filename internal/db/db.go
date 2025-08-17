package db

import (
	"fmt"
	"time"

	"github.com/aigic8/corn/internal/db/models"
	"github.com/aigic8/corn/internal/db/schema"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DbClient = gorm.DB

var ErrNotExist = gorm.ErrRecordNotFound

type (
	Db struct {
		DbAddr  string
		Timeout time.Duration
		c       *DbClient
		Run     *models.RunModel
	}
)

func NewDb(dbAddr string, defaultTimeout time.Duration) (*Db, error) {
	c, err := gorm.Open(sqlite.Open(dbAddr), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("connecting to db: %w", err)
	}

	runModel := models.NewRunModel(c, defaultTimeout)

	return &Db{DbAddr: dbAddr, c: c, Run: runModel, Timeout: defaultTimeout}, nil
}

func (db *Db) Init() error {
	return db.c.AutoMigrate(&schema.Run{})
}
