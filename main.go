package main

import (
	"fmt"

	"github.com/aigic8/corn/lib/config"
	"github.com/aigic8/corn/lib/logs"
	"github.com/aigic8/corn/lib/runner"
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

	logger, err := logs.NewLogger(c.LogsDir)
	if err != nil {
		panic(fmt.Errorf("creating a new logger: %w", err))
	}

	r, err := runner.NewRunner(c, logger)
	if err != nil {
		panic(fmt.Errorf("failed to create a new runner: %w", err))
	}

	if err := r.ScheduleJobs(); err != nil {
		panic(fmt.Errorf("failed to schedule jobs: %w", err))
	}
}
