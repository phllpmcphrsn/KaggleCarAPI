package main

import (
	"os"
	"path/filepath"
	"strings"

	log "golang.org/x/exp/slog"

	"github.com/spf13/viper"
)

// Config holds the configuration values
type Config struct {
	Environment string
	API         APIConfig
	Log         LogLevel
	Database    DatabaseConfig
	CSV         CSVConfig
}

// APIConfig holds the API configuration values
type APIConfig struct {
	Address string
	Path    string
}

// LogLevel holds the log configuration values
type LogLevel struct {
	LevelStr string
	Level    log.Level
}

// DatabaseConfig holds the database configuration values
type DatabaseConfig struct {
	Host string
	Port int
	Name string
	SSL  SSL
}

// SSL determines if SSL will be enabled for database connections
type SSL struct {
	Enabled bool
}

// CSVConfig holds the CSV configuration values.
type CSVConfig struct {
	Filename string
}

// LoadConfig loads the configuration values from the specified file.
func LoadConfig(file string) (*Config, error) {
	// Set the file name and path
	if file != "" {
		viper.SetConfigFile(file)
	} else {
		// If no file specified, look for the default file in the current directory
		dir, err := os.Getwd()
		if err != nil {
			log.Error("Failed to get current directory", "err", err, "")
		}
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(dir)
	}

	// Read the configuration file
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	// Unmarshal the configuration values into the Config struct
	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}
	log.Info(config.Log.LevelStr)
	// Get a valid slog log level
	config.Log.Level = GetLogLevel(config.Log.LevelStr)

	return &config, nil
}

// GetConfigFilePath returns the absolute path of the config file based on the current directory.
func GetConfigFilePath() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Error("Failed to get current directory", "err", err)
	}
	return filepath.Join(dir, "config.yaml")
}

// GetLogLevel returns the slog log level based on a string representation of the log level.
// INFO is used as the default in case a level isn't given or is unexpected
func GetLogLevel(logLevel string) log.Level {
	level := strings.ToLower(logLevel)
	print(level)
	log.Info(level)
	switch level {
	case "debug":
		return log.LevelDebug
	case "info":
		return log.LevelInfo
	case "warn":
		return log.LevelWarn
	case "error":
		return log.LevelError
	default:
		log.Warn("Supported log levels are: debug, info, warn, and error")
		return log.LevelInfo
	}
}
