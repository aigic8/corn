package runner

import (
	"fmt"
	"strings"

	"github.com/aigic8/corn/lib/common"
	"github.com/aigic8/corn/lib/config"
	"github.com/aigic8/corn/lib/logs"
	"github.com/go-co-op/gocron/v2"
)

type (
	Runner struct {
		L         *logs.Logger
		Config    *config.Config
		Scheduler gocron.Scheduler
	}

	StringWriter interface {
		Write(p []byte) (int, error)
		String() string
	}
)

func NewRunner(c *config.Config, l *logs.Logger) (*Runner, error) {
	s, err := gocron.NewScheduler()
	if err != nil {
		return nil, fmt.Errorf("creating a new scheduler: %w", err)
	}
	return &Runner{L: l, Config: c, Scheduler: s}, nil
}

func (r *Runner) ScheduleJobs() error {
	for jobName, job := range r.Config.Jobs {
		for _, schedule := range job.Schedules {
			_, err := r.Scheduler.NewJob(gocron.CronJob(schedule, true), gocron.NewTask(r.JobFunc(jobName)))
			if err != nil {
				return fmt.Errorf("scheduling job '%s' with schedule '%s': %w", jobName, schedule, err)
			}
			r.JobFunc(jobName)()
		}
	}
	r.L.L.Debug().Msg("scheduling jobs were successful")
	return nil
}

func (r *Runner) JobFunc(jobName string) func() {
	return func() {
		job := r.Config.Jobs[jobName]

		jobLogger, closeJobLogger, err := r.L.NewJobLogger(jobName)
		if err != nil {
			r.L.L.Err(fmt.Errorf("creating job logger for job '%s': %w", jobName, err)).Msg("failed to create job logger")
			return
		}
		defer closeJobLogger()
		cmd, args := common.SeparateCommandFromArgs(job.Command)

		stdoutWriter := &strings.Builder{}
		var stderrWriter StringWriter
		if !job.IgnoreStdErrLog {
			stderrWriter = &strings.Builder{}
		} else {
			stderrWriter = &common.StringNullWriter{}
		}

		failed := false
		err = common.RunCommand(&common.RunCommandOpts{
			Cmd:    cmd,
			Args:   args,
			Stdout: stdoutWriter,
			Stderr: stderrWriter,
		})
		if err != nil {
			failed = true
			// TODO: notify
			r.L.L.Err(fmt.Errorf("running job '%s': %w", jobName, err)).Msg("Running job failed")
		}
		if !job.OnlyLogOnFail || (job.OnlyLogOnFail && failed) {
			log := jobLogger.Err(err).Str("stdout", stdoutWriter.String()).Bool("failed", failed)
			if !job.IgnoreStdErrLog {
				log.Str("stderr", stderrWriter.String())
			}
			log.Msgf("job '%s' executed", jobName)
		}
	}
}
