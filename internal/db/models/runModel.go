package models

import (
	"context"
	"fmt"

	"github.com/aigic8/corn/internal/db/schema"
	"gorm.io/gorm"
)

type (
	RunModel struct {
		c *gorm.DB
	}
)

func NewRunModel(c *gorm.DB) *RunModel {
	return &RunModel{c: c}
}

func (rm *RunModel) Create(job string) (uint, error) {
	// TODO: add timeout
	run := schema.Run{Job: job}
	if err := gorm.G[schema.Run](rm.c, gorm.WithResult()).Create(context.Background(), &run); err != nil {
		return 0, fmt.Errorf("creating run item: %w", err)
	}
	return run.ID, nil
}

func (rm *RunModel) Get(id uint) (*schema.Run, error) {
	// TODO: add timeout
	runData, err := gorm.G[schema.Run](rm.c).Where("id = ?", id).First(context.Background())
	if err != nil {
		return nil, fmt.Errorf("getting retry counts: %w", err)
	}
	return &runData, nil
}

func (rm *RunModel) UpdateRunRetries(id uint, retries uint) error {
	// TODO: add timeout
	_, err := gorm.G[schema.Run](rm.c).Where("id = ?", id).Update(context.Background(), "retries", retries)
	if err != nil {
		return fmt.Errorf("updating retries: %w", err)
	}
	return nil
}
