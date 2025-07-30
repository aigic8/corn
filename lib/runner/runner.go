package runner

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aigic8/corn/internal/db"
	"github.com/aigic8/corn/lib/command"
	"github.com/aigic8/corn/lib/common"
	"github.com/aigic8/corn/lib/config"
	"github.com/aigic8/corn/lib/logs"
	"github.com/aigic8/corn/lib/notif"
	"github.com/go-co-op/gocron/v2"
	"github.com/nikoksr/notify"
)

type FailureNotifier = func(err error, logMsg string, stdout string, stderr string)
type (
	Runner struct {
		L         *logs.Logger
		Config    *config.Config
		Scheduler gocron.Scheduler
		Notif     *notif.Notif
		p         *command.CommandParser
		db        *db.Db
	}

	StringWriter interface {
		Write(p []byte) (int, error)
		String() string
	}
)

func NewRunner(c *config.Config, l *logs.Logger, db *db.Db) (*Runner, error) {
	s, err := gocron.NewScheduler()
	if err != nil {
		return nil, fmt.Errorf("creating a new scheduler: %w", err)
	}

	notifiers, err := CompileNotifServices(c.Notifiers, c.DisableNotifications)
	if err != nil {
		return nil, fmt.Errorf("compiling notifiers: %w", err)
	}
	n := notif.NewNotif(time.Duration(c.NotifyTimeoutMs)*time.Millisecond, notifiers, c.DisableNotifications)

	return &Runner{L: l, Config: c, Scheduler: s, Notif: n, p: command.NewCommandParser(), db: db}, nil
}

func (r *Runner) ScheduleJobs() error {
	for jobName, job := range r.Config.Jobs {
		for _, schedule := range job.Schedules {
			_, err := r.Scheduler.NewJob(gocron.CronJob(schedule, true), gocron.NewTask(r.JobFunc(jobName, false)))
			if err != nil {
				return fmt.Errorf("scheduling job '%s' with schedule '%s': %w", jobName, schedule, err)
			}
		}
	}
	r.L.L.Debug().Msg("scheduling jobs were successful")
	return nil
}

func (r *Runner) JobFunc(jobName string, test bool) func() {
	failStrategy := r.Config.Jobs[jobName].FailStrategy
	if !test && failStrategy.Retry != nil {
		return r.retryJobFunc(jobName, failStrategy.Retry)
	} else if !test && failStrategy.Halt != nil {
		r.L.L.Err(fmt.Errorf("fail strategy halt is not yet implemented (job '%s')", jobName)).Msg("fail strategy halt not implemented")
		return r.normalJobFunc(jobName)
	} else {
		// default fail strategy is continue
		return r.normalJobFunc(jobName)
	}
}

func (r *Runner) normalJobFunc(jobName string) func() {
	return func() {
		jobLogger, closeJobLogger, err := r.L.NewJobLogger(jobName)
		if err != nil {
			err := fmt.Errorf("creating job logger for job '%s': %w", jobName, err)
			r.L.L.Err(err).Msg("failed to create job logger")
		}
		defer closeJobLogger()

		notifierName := r.getNotifierForJob(jobName, true)
		notifFailure := r.failureLogNotifier(jobName, &jobLogger, notifierName)
		r.runJob(jobName, notifFailure, &jobLogger)
	}
}

func (r *Runner) retryJobFunc(jobName string, retry *config.FailStrategyRetry) func() {
	return func() {
		jobLogger, closeJobLogger, err := r.L.NewJobLogger(jobName)
		if err != nil {
			err := fmt.Errorf("creating job logger for job '%s': %w", jobName, err)
			r.L.L.Err(err).Msg("failed to create job logger")
		}
		defer closeJobLogger()

		notifierName := r.getNotifierForJob(jobName, true)
		notifFailure := r.failureLogNotifier(jobName, &jobLogger, notifierName)

		var retriesCount uint = 0
		retriesCount, err = r.db.Retry.GetRetryCount(jobName)
		if err != nil {
			notifFailure(fmt.Errorf("getting number of retries for job '%s': %w", jobName, err), "could not get number of retries", "", "")
			return
		}

		if retriesCount+1 >= retry.MaxRetries {
			if err = r.db.Retry.UpsertRetries(jobName, 0); err != nil {
				notifFailure(fmt.Errorf("upserting retry count for job '%s': %w", jobName, err), "could not set retry count", "", "")
				return
			}
			return
		} else {
			if err = r.db.Retry.UpsertRetries(jobName, retriesCount+1); err != nil {
				notifFailure(fmt.Errorf("upserting retry count for job '%s': %w", jobName, err), "could not set retry count", "", "")
				return
			}
			if err = r.runJob(jobName, notifFailure, &jobLogger); err != nil {
				nextTime := gocron.OneTimeJobStartImmediately()
				if retry.CoolOffSecs != 0 {
					timeVal := time.Now().Add(time.Duration(retry.CoolOffSecs) * time.Second)
					nextTime = gocron.OneTimeJobStartDateTime(timeVal)
				}

				r.Scheduler.NewJob(
					gocron.OneTimeJob(nextTime),
					gocron.NewTask(r.JobFunc(jobName, false)),
				)
			} else {
				if err = r.db.Retry.UpsertRetries(jobName, 0); err != nil {
					notifFailure(fmt.Errorf("upserting retry count for job '%s': %w", jobName, err), "could not set retry count", "", "")
					return
				}
			}
		}
	}

}

func (r *Runner) runJob(jobName string, notifFailure FailureNotifier, jobLogger *logs.InternalLogger) error {
	job, exists := r.Config.Jobs[jobName]
	if !exists {
		err := fmt.Errorf("job with name '%s' does not exist", jobName)
		r.L.L.Err(err).Msg("failed to find the job")
		return err
	}

	var stdout, stderr string
	var err error
	if job.RemoteName == "" {
		err, stdout, stderr = r.runCommandLocally(jobName)
	} else {
		err, stdout, stderr = r.runCommandRemotely(jobName)
	}
	if err != nil {
		notifFailure(fmt.Errorf("job '%s' failed: %w", jobName, err), fmt.Sprintf("job '%s' failed", jobName), stdout, stderr)
		return err
	} else {
		if !job.OnlyLogOnFail {
			jobLogger.Info().Bool("failed", false).Str("stdout", stdout).Msgf("job '%s' executed", jobName)
		}
		if !job.OnlyNotifyOnFail {
			send := r.Notif.UseService(r.getNotifierForJob(jobName, false))
			send(fmt.Sprintf("job '%s' executed", jobName), "Job executed successfully")
			if err != nil {
				r.L.L.Err(fmt.Errorf("sending notification: %w", jobName, err)).Msg("sending notification failed")
			}
		}
	}

	return nil
}

// logs the error and notifies the user on failure of the job
func (r *Runner) failureLogNotifier(jobName string, log *logs.InternalLogger, notifServiceName string) FailureNotifier {
	if notifServiceName != "" {
		sendNotification := r.Notif.UseService(notifServiceName)
		return func(err error, logMsg, stdout, stderr string) {
			log := log.Err(err).Str("stdout", stdout)
			if !r.Config.Jobs[jobName].IgnoreStderrLog {
				log.Str("stderr", stderr)
			}
			log.Msg(logMsg)
			if err = sendNotification(fmt.Sprintf("job '%s' failed", jobName), fmt.Sprintf("error: %s", err.Error())); err != nil {
				r.L.L.Err(err).Msg("failed sending notification for failure")
			}
		}
	} else {
		return func(err error, logMsg, stdout, stderr string) {}
	}
}

func (r *Runner) getNotifierForJob(jobName string, failure bool) string {
	if r.Config.DisableNotifications {
		return ""
	}
	job := r.Config.Jobs[jobName]
	if failure {
		if job.FailNotifier != "" {
			return job.FailNotifier
		} else if r.Config.DefaultFailNotifier != "" {
			return r.Config.DefaultFailNotifier
		}
	}
	if job.Notifier != "" {
		return job.Notifier
	} else if r.Config.DefaultNotifier != "" {
		return r.Config.DefaultNotifier
	}
	return ""
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

func (r *Runner) runCommandLocally(jobName string) (error, string, string) {
	cmd := r.Config.Jobs[jobName].Command
	parsed, err := r.p.Parse(strings.NewReader(cmd))
	if err != nil {
		return fmt.Errorf("parsing command: %w", jobName, err), "", ""
	}
	stdoutWriter := &strings.Builder{}
	var stderrWriter StringWriter
	if !r.Config.Jobs[jobName].IgnoreStderrLog {
		stderrWriter = &strings.Builder{}
	} else {
		stderrWriter = &common.StringNullWriter{}
	}

	err = command.RunCommand(&command.RunCommandOpts{
		Cmd:     parsed,
		Stdout:  stdoutWriter,
		Stderr:  stderrWriter,
		Timeout: r.getTimeoutForJob(jobName),
	})
	if err != nil {
		return fmt.Errorf("running command: %w", jobName, err), stdoutWriter.String(), stderrWriter.String()
	}
	return nil, stdoutWriter.String(), stderrWriter.String()
}

func (r *Runner) runCommandRemotely(jobName string) (error, string, string) {
	remoteName, remote := r.getRemoteForJob(jobName)
	if remote == nil {
		return fmt.Errorf("remote '%s' for job '%s' was not found", remoteName, jobName), "", ""
	}

	client, err := command.LoginToRemote(remote)
	if err != nil {
		return fmt.Errorf("logging in to remote '%s': %w", remoteName, err), "", ""
	}

	timeout := r.getTimeoutForJob(jobName)
	cmd := r.Config.Jobs[jobName].Command
	// TODO: seperate the stderr and stdout
	var combinedOutput []byte
	if timeout == 0 {
		combinedOutput, err = client.Run(cmd)
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		combinedOutput, err = client.RunContext(ctx, cmd)
	}
	if err != nil {
		return fmt.Errorf("running command: %w", err), string(combinedOutput), ""
	}
	return nil, string(combinedOutput), ""
}

func (r *Runner) getRemoteForJob(jobName string) (string, *config.Remote) {
	remoteName := r.Config.Jobs[jobName].RemoteName
	remote := r.Config.Remotes[remoteName]
	return remoteName, &remote
}

// helper function to convert config notification services into application notifiers
func CompileNotifServices(services map[string]config.NotifyService, disabled bool) (map[string][]notify.Notifier, error) {
	if disabled {
		return map[string][]notify.Notifier{}, nil
	}
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
