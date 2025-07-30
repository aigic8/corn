package config

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/aigic8/corn/lib/common"
	"github.com/go-playground/validator/v10"
)

// sqlite db file path, home path should be added later
const DEFAULT_DB_PATH = ".corn/corn.sqlite"

// logs dir, home path should be added later
const DEFAULT_LOGS_DIR = ".corn/logs"
const DEFAULT_NOTIFY_TIMEOUT = 10000

type (
	Config struct {
		DbAddr               string                   `yaml:"dbAddr"`
		Jobs                 map[string]Job           `yaml:"jobs" validate:"required,min=1,dive"`
		Remotes              map[string]Remote        `yaml:"remotes" validate:"dive"`
		Notifiers            map[string]NotifyService `yaml:"notifiers" validate:"dive"`
		LogsDir              string                   `yaml:"logsDir"`
		NotifyTimeoutMs      int                      `yaml:"notifyTimeoutMs"`
		DefaultFailNotifier  string                   `yaml:"defaultFailNotifier"`
		DefaultNotifier      string                   `yaml:"defaultNotifier"`
		DefaultTimeoutS      int                      `yaml:"defaultTimeoutS"`
		DisableNotifications bool                     `yaml:"disableNotifications"`
		Timezone             string                   `yaml:"timezone"`
	}

	NotifyService struct {
		Telegram []TelegramNotifyService `yaml:"telegram"`
		Discord  []DiscordNotifyService  `yaml:"discord"`
	}

	Remote struct {
		Username string      `yaml:"username" validate:"required,min=1"`
		Address  string      `yaml:"address" validate:"required,min=1"`
		Port     uint        `yaml:"port" validate:"required"`
		Auth     *RemoteAuth `yaml:"auth" validate:"required"`
	}

	RemoteAuth struct {
		PasswordAuth *PasswordAuth `yaml:"passwordAuth"`
		KeyAuth      *KeyAuth      `yaml:"keyAuth"`
	}

	PasswordAuth struct {
		Password string `yaml:"password" validate:"required"`
	}

	KeyAuth struct {
		KeyPath    string `yaml:"keyPath" validate:"required"`
		Passphrase string `yaml:"passphrase"`
	}

	TelegramNotifyService struct {
		Token     string  `yaml:"token" validate:"required"`
		Receivers []int64 `yaml:"receivers" validate:"required,min=1"`
	}

	DiscordNotifyService struct {
		OAuth2Token string   `yaml:"oAuth2Token"`
		BotToken    string   `yaml:"botToken"`
		Channels    []string `yaml:"channels" validate:"required,min=1"`
	}

	Job struct {
		Schedules        []string      `yaml:"schedules" validate:"required"`
		Command          string        `yaml:"command" validate:"required"`
		OnlyLogOnFail    bool          `yaml:"onlyLogOnFail"`
		IgnoreStderrLog  bool          `yaml:"ignoreStdErrLog"`
		OnlyNotifyOnFail bool          `yaml:"onlyNotifyOnFail"`
		FailNotifier     string        `yaml:"failNotifier"`
		Notifier         string        `yaml:"notifier"`
		TimeoutS         int           `yaml:"timeoutS"`
		RemoteName       string        `yaml:"remoteName"`
		FailStrategy     *FailStrategy `yaml:"failStrategy"`
	}

	FailStrategy struct {
		Retry *FailStrategyRetry `yaml:"retry"`
		Halt  *FailStrategyHalt  `yaml:"halt"`
	}

	FailStrategyRetry struct {
		MaxRetries  uint `yaml:"maxRetries" validate:"required"`
		CoolOffSecs uint `yaml:"coolOffSecs"`
	}

	FailStrategyHalt struct {
	}
)

var CONFIG_PATHS = []string{
	".config/corn/corn.yaml",
	".config/corn/corn.yml",
}

func ParseConfig(configPath string) (*Config, error) {
	configBytes, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(configBytes, &config); err != nil {
		return nil, fmt.Errorf("parsing config file: %w", err)
	}

	// set the default values if not exist
	if config.LogsDir == "" {
		homePaths, err := common.FromHome(DEFAULT_LOGS_DIR)
		if err != nil {
			return nil, fmt.Errorf("getting home directory: %w", err)
		}
		config.LogsDir = homePaths[0]
	}
	if config.DbAddr == "" {
		homePaths, err := common.FromHome(DEFAULT_DB_PATH)
		if err != nil {
			return nil, fmt.Errorf("getting home directory: %w", err)
		}
		config.DbAddr = homePaths[0]
	}
	if config.NotifyTimeoutMs == 0 {
		config.NotifyTimeoutMs = DEFAULT_NOTIFY_TIMEOUT
	}

	return &config, nil
}

func ParseAndValidateConfig(configPath string, v *validator.Validate) (*Config, error) {
	config, err := ParseConfig(configPath)
	if err != nil {
		return nil, err
	}
	if err := v.Struct(config); err != nil {
		return config, fmt.Errorf("validating config: %w", err)
	}
	return config, nil
}

func GetConfigPath() (string, error) {
	configPaths, err := common.FromHome(CONFIG_PATHS...)
	if err != nil {
		return "", fmt.Errorf("appending home path to config paths: %w", err)
	}
	for _, configPath := range configPaths {
		stat, err := os.Stat(configPath)
		if err != nil || stat.IsDir() {
			continue
		} else {
			return configPath, nil
		}
	}
	return "", fmt.Errorf("no config file found. Checked paths:\n%s", strings.Join(configPaths, "\n"))
}
