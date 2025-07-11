package runner

import (
	"fmt"
	"strings"
	"time"

	"github.com/aigic8/corn/lib/command"
	"github.com/aigic8/corn/lib/common"
	"github.com/aigic8/corn/lib/config"
	"github.com/aigic8/corn/lib/logs"
	"github.com/aigic8/corn/lib/notif"
	"github.com/go-co-op/gocron/v2"
	"github.com/nikoksr/notify"
)

type (
	Runner struct {
		L         *logs.Logger
		Config    *config.Config
		Scheduler gocron.Scheduler
		Notif     *notif.Notif
		p         *command.CommandParser
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

	notifiers, err := CompileNotifServices(c.Notifiers)
	if err != nil {
		return nil, fmt.Errorf("compiling notifiers: %w", err)
	}
	n := notif.NewNotif(time.Duration(c.NotifyTimeoutMs)*time.Millisecond, notifiers, c.DisableNotifications)

	return &Runner{L: l, Config: c, Scheduler: s, Notif: n, p: command.NewCommandParser()}, nil
}

func (r *Runner) ScheduleJobs() error {
	for jobName, job := range r.Config.Jobs {
		for _, schedule := range job.Schedules {
			_, err := r.Scheduler.NewJob(gocron.CronJob(schedule, true), gocron.NewTask(r.JobFunc(jobName)))
			if err != nil {
				return fmt.Errorf("scheduling job '%s' with schedule '%s': %w", jobName, schedule, err)
			}
		}
	}
	r.L.L.Debug().Msg("scheduling jobs were successful")
	return nil
}

func (r *Runner) JobFunc(jobName string) func() {
	return func() {
		job, exists := r.Config.Jobs[jobName]
		if !exists {
			r.L.L.Err(fmt.Errorf("job with name '%s' does not exist", jobName)).Msg("failed to find the job")
			return
		}

		jobLogger, closeJobLogger, err := r.L.NewJobLogger(jobName)
		if err != nil {
			r.L.L.Err(fmt.Errorf("creating job logger for job '%s': %w", jobName, err)).Msg("failed to create job logger")
			return
		}
		defer closeJobLogger()

		parsed, err := r.p.Parse(strings.NewReader(job.Command))
		if err != nil {
			r.L.L.Err(fmt.Errorf("parsing command for job '%s': %w", jobName, err)).Msg("failed to parse job")
			return
		}

		// FIXME: notify the user if the command fail before running (remove all the logging and returns before this comment)

		stdoutWriter := &strings.Builder{}
		var stderrWriter StringWriter
		if !job.IgnoreStderrLog {
			stderrWriter = &strings.Builder{}
		} else {
			stderrWriter = &common.StringNullWriter{}
		}

		failed := false
		err = command.RunCommand(&command.RunCommandOpts{
			Cmd:     parsed,
			Stdout:  stdoutWriter,
			Stderr:  stderrWriter,
			Timeout: r.getTimeoutForJob(jobName),
		})
		if err != nil {
			failed = true
			r.L.L.Err(fmt.Errorf("running job '%s': %w", jobName, err)).Msg("running job failed")
		}

		// handle logging
		if !job.OnlyLogOnFail || (job.OnlyLogOnFail && failed) {
			log := jobLogger.Err(err).Str("stdout", stdoutWriter.String()).Bool("failed", failed)
			if !job.IgnoreStderrLog {
				log.Str("stderr", stderrWriter.String())
			}
			log.Msgf("job '%s' executed", jobName)
		}

		// handle notification
		if !job.OnlyNotifyOnFail || (job.OnlyNotifyOnFail && failed) {
			// FIXME: do not error and send notification if the user has notifications disabled
			notifierName, notifierErr := r.getNotifierForJob(jobName, true)
			if notifierErr != nil {
				r.L.L.Err(fmt.Errorf("getting notifier for job '%s': %w", jobName, err)).Msg("getting notifier failed")
				return
			}
			send := r.Notif.UseService(notifierName)
			if failed {
				err = send(fmt.Sprintf("job '%s' failed", jobName), fmt.Sprintf("error: %s\nstdout: %s", err.Error(), stdoutWriter.String()))
				if err != nil {
					r.L.L.Err(fmt.Errorf("sending notification: %w", jobName, err)).Msg("sending notification failed")
				}
			} else {
				send(fmt.Sprintf("job '%s' executed", jobName), "Job executed successfully")
				if err != nil {
					r.L.L.Err(fmt.Errorf("sending notification: %w", jobName, err)).Msg("sending notification failed")
				}
			}
		}
	}
}

func (r *Runner) getNotifierForJob(jobName string, failure bool) (string, error) {
	job := r.Config.Jobs[jobName]
	if failure {
		if job.FailNotifier != "" {
			return job.FailNotifier, nil
		} else if r.Config.DefaultFailNotifier != "" {
			return r.Config.DefaultFailNotifier, nil
		}
	}
	if job.Notifier != "" {
		return job.Notifier, nil
	} else if r.Config.DefaultNotifier != "" {
		return r.Config.DefaultNotifier, nil
	}

	return "", fmt.Errorf("no notifier found for job '%s'", jobName)
}

// returns the timeout for the job based on the default timeout
// and the job's timeout.
// If not timeout was found, it would return 0 time duration
func (r *Runner) getTimeoutForJob(jobName string) time.Duration {
	jobTimeout := r.Config.Jobs[jobName].TimeoutS
	if jobTimeout != 0 {
		return time.Duration(jobTimeout) * time.Second
	}
	if r.Config.DefaultTimeoutS != 0 {
		return time.Duration(r.Config.DefaultTimeoutS) * time.Second
	}
	return 0
}

// helper function to convert config notification services into application notifiers
func CompileNotifServices(services map[string]config.NotifyService) (map[string][]notify.Notifier, error) {
	res := make(map[string][]notify.Notifier, len(services))
	for serviceName, service := range services {
		res[serviceName] = []notify.Notifier{}
		if service.Telegram != nil {
			for _, telegramApp := range service.Telegram {
				notifier, err := notif.NewTelegramNotifier(telegramApp.Token, telegramApp.Receivers...)
				if err != nil {
					return res, fmt.Errorf("creating telegram notifier for service '%s': %w", serviceName, err)
				}
				res[serviceName] = append(res[serviceName], notifier)
			}
		}
		if service.Discord != nil {
			for _, discordApp := range service.Discord {
				notifier, err := notif.NewDiscordNotifier(&notif.DiscordNotifierOpts{
					BotToken:    discordApp.BotToken,
					OAuth2Token: discordApp.OAuth2Token,
					ChanelIDs:   discordApp.Channels,
				})
				if err != nil {
					return res, fmt.Errorf("creating discord notifier for service '%s': %w", serviceName, err)
				}
				res[serviceName] = append(res[serviceName], notifier)
			}
		}
	}
	return res, nil
}
