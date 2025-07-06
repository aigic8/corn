package config

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/aigic8/corn/lib/common"
	"github.com/go-playground/validator/v10"
)

const DEFAULT_LOGS_DIR = ".corn/logs"

var validate *validator.Validate

type (
	Config struct {
		Jobs    map[string]Job `yaml:"jobs" validate:"required,min=1"`
		LogsDir string         `yaml:"logsDir"`
	}

	Job struct {
		Schedules       []string `yaml:"schedules" validate:"required"`
		Command         string   `yaml:"command" validate:"required"`
		OnlyLogOnFail   bool     `yaml:"onlyLogOnFail"`
		IgnoreStdErrLog bool     `yaml:"ignoreStdErrLog"`
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

	// set the default logs dir if not exists
	if config.LogsDir == "" {
		homePaths, err := common.FromHome(DEFAULT_LOGS_DIR)
		if err != nil {
			return nil, fmt.Errorf("getting home directory: %w", err)
		}
		config.LogsDir = homePaths[0]
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
