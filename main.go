package main

import (
	"fmt"

	"github.com/aigic8/corn/lib/common"
	"github.com/aigic8/corn/lib/config"
	"github.com/go-co-op/gocron/v2"
	"github.com/go-playground/validator/v10"
)

func main() {
	configPath, err := config.GetConfigPath()
	if err != nil {
		panic(err)
	}

	v := validator.New(validator.WithRequiredStructEnabled())
	c, err := config.ParseAndValidateConfig(configPath, v)
	if err != nil {
		panic(err)
	}

	s, err := gocron.NewScheduler()
	if err != nil {
		panic(fmt.Errorf("creating a new scheduler: %w", err))
	}

	for _, job := range c.Jobs {
		for _, schedule := range job.Schedules {
			s.NewJob(gocron.CronJob(schedule, true), gocron.NewTask(func() {
				cmd, args := common.SeparateCommandFromArgs(job.Command)
				if err := common.RunCommand(cmd, args); err != nil {
					// TODO: error handling
					panic(err)
				}
			}))
		}
	}
}
