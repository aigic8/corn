package models

import (
	"fmt"

	"github.com/aigic8/corn/internal/db/schema"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type (
	RetryModel struct {
		c *gorm.DB
	}
)

func NewRetryModel(c *gorm.DB) *RetryModel {
	return &RetryModel{c: c}
}

func (rm *RetryModel) GetRetryCount(job string) (uint, error) {
	var retries uint = 0
	if err := rm.c.Select("retries").Where(&schema.Retry{Job: job}).Scan(retries); err != nil {
		return 0, fmt.Errorf("getting retry counts: %w", err)
	}
	return retries, nil
}

func (rm *RetryModel) UpsertRetries(job string, retries uint) error {
	row := schema.Retry{Job: job, Retries: retries}
	return rm.c.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "job"}},
		DoUpdates: clause.AssignmentColumns([]string{"retries"}),
	}).Create(&row).Error
}
