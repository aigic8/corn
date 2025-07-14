package logs

import (
	"fmt"
	"os"
	"path"

	"github.com/aigic8/corn/lib/common"
	"github.com/rs/zerolog"
)

type InternalLogger = zerolog.Logger

type (
	Logger struct {
		L       InternalLogger
		logFile *os.File
		LogsDir string
	}
)

const LOG_FILE_NAME = "corn.jsonl"
const LOG_DIR_PERM = 0750
const LOG_FILE_PERM = 0640

func NewLogger(logBaseDir string) (*Logger, error) {
	if err := common.MakeDirAllIfNotExist(logBaseDir, LOG_DIR_PERM); err != nil {
		return nil, fmt.Errorf("creating log base directory (%s): %w", logBaseDir, err)
	}

	logFilePath := path.Join(logBaseDir, LOG_FILE_NAME)
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, LOG_FILE_PERM)
	if err != nil {
		return nil, fmt.Errorf("creating log file (%s): %w", logFilePath, err)
	}

	return &Logger{L: zerolog.New(logFile).With().Timestamp().Logger(), LogsDir: logBaseDir}, nil
}

// creates a logger for a job, returns the logger, close function and error
// the close function should be used (probably deferred) to close the job file stream
// IMPORTANT: the close function might be null if there is an error
func (l *Logger) NewJobLogger(jobName string) (zerolog.Logger, func() error, error) {
	logsFilePath := path.Join(l.LogsDir, "/jobs/", jobName+".jsonl")
	baseDir := path.Dir(logsFilePath)
	if err := common.MakeDirAllIfNotExist(baseDir, LOG_DIR_PERM); err != nil {
		return l.L, nil, fmt.Errorf("creating jobs base directory (%s): %w", baseDir, err)
	}

	logsFile, err := os.OpenFile(logsFilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, LOG_FILE_PERM)
	close := func() error {
		return logsFile.Close()
	}
	if err != nil {
		return l.L, close, fmt.Errorf("creating jobs log file: %w", err)
	}

	return zerolog.New(logsFile).With().Timestamp().Str("job", jobName).Logger(), close, nil
}

func (l *Logger) Close() error {
	return l.logFile.Close()
}

// TODO: create a function called clean which removes the logs older than specific time
// and run it periodically
