package models

import (
	"context"
	"fmt"
	"time"

	"github.com/aigic8/corn/internal/db/schema"
	"gorm.io/gorm"
)

type (
	RunModel struct {
		c       *gorm.DB
		Timeout time.Duration
	}
)

func NewRunModel(c *gorm.DB, defaultTimeout time.Duration) *RunModel {
	return &RunModel{c: c, Timeout: defaultTimeout}
}

func (rm *RunModel) Create(job string) (uint, error) {
	ctx, cancel := context.WithTimeout(context.Background(), rm.Timeout)
	defer cancel()
	run := schema.Run{Job: job}
	if err := gorm.G[schema.Run](rm.c, gorm.WithResult()).Create(ctx, &run); err != nil {
		return 0, fmt.Errorf("creating run item: %w", err)
	}
	return run.ID, nil
}

func (rm *RunModel) Get(id uint) (*schema.Run, error) {
	ctx, cancel := context.WithTimeout(context.Background(), rm.Timeout)
	defer cancel()
	runData, err := gorm.G[schema.Run](rm.c).Where("id = ?", id).First(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting retry counts: %w", err)
	}
	return &runData, nil
}

func (rm *RunModel) UpdateRunRetries(id uint, retries uint) error {
	ctx, cancel := context.WithTimeout(context.Background(), rm.Timeout)
	defer cancel()
	_, err := gorm.G[schema.Run](rm.c).Where("id = ?", id).Update(ctx, "retries", retries)
	if err != nil {
		return fmt.Errorf("updating retries: %w", err)
	}
	return nil
}
