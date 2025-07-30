package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"
	_ "time/tzdata"

	"github.com/aigic8/corn/internal/db"
	"github.com/aigic8/corn/lib/common"
	"github.com/aigic8/corn/lib/config"
	"github.com/aigic8/corn/lib/logs"
	"github.com/aigic8/corn/lib/runner"
	"github.com/alexflint/go-arg"
	"github.com/go-playground/validator/v10"
)

type (
	Args struct {
		Test *TestCommand `arg:"subcommand" help:"test a job"`
	}

	TestCommand struct {
		Job string `arg:"-j,--job,required" help:"name of the job to test"`
	}
)

func main() {
	var args Args
	arg.MustParse(&args)

	configPath, err := config.GetConfigPath()
	if err != nil {
		panic(err)
	}

	v := validator.New(validator.WithRequiredStructEnabled())
	c, err := config.ParseAndValidateConfig(configPath, v)
	if err != nil {
		panic(fmt.Sprintf("parsing and validating config from '%s': %s", configPath, err.Error()))
	}

	logger, err := logs.NewLogger(c.LogsDir)
	if err != nil {
		panic(fmt.Sprintf("creating a new logger: %s", err.Error()))
	}

	logger.L.Debug().Msgf("loaded config file from '%s'", configPath)

	// source: https://stackoverflow.com/a/64769139
	if c.Timezone != "" {
		loc, err := time.LoadLocation(c.Timezone)
		if err != nil {
			panic(fmt.Sprintf("setting the timezone: %s", err.Error()))
		}
		time.Local = loc
		logger.L.Debug().Msgf("set timezone to '%s'", c.Timezone)
	}

	// TODO: check if the path of the db is a file
	dbDir := path.Dir(c.DbAddr)
	if err = common.MakeDirAllIfNotExist(dbDir, 0750); err != nil {
		panic(fmt.Sprintf("creating directory '%s' for database: %w", err))
	}
	logger.L.Debug().Msgf("created directory '%s' for db", dbDir)

	db, err := db.NewDb(c.DbAddr)
	if err != nil {
		panic(fmt.Sprintf("creating new db with address '%s': %w", c.DbAddr, err))
	}

	r, err := runner.NewRunner(c, logger, db)
	if err != nil {
		panic(fmt.Errorf("failed to create a new runner: %w", err))
	}

	if args.Test != nil {
		r.JobFunc(args.Test.Job, true)()
	} else {
		if err := r.ScheduleJobs(); err != nil {
			panic(fmt.Errorf("failed to schedule jobs: %w", err))
		}

		if err := r.RunJobs(); err != nil {
			panic(fmt.Errorf("failed to run jobs: %w", err))
		}
		defer r.Shutdown()

		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		fmt.Println("Running... Press Ctrl-C to stop.")

		<-ctx.Done()
	}

}
