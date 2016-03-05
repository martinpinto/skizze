package config

import (
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/njpatel/loggo"

	"utils"
)

var logger = loggo.GetLogger("config")

// Config stores all configuration parameters for Go
type Config struct {
	InfoDir              string `toml:"info_dir"`
	DataDir              string `toml:"data_dir"`
	Port                 uint   `toml:"port"`
	SaveThresholdSeconds uint   `toml:"save_threshold_seconds"`
}

var config *Config

// MaxKeySize ...
const MaxKeySize int = 32768 // max key size BoltDB in bytes

func parseConfigTOML() *Config {
	configPath := os.Getenv("SKZ_CONFIG")
	if configPath == "" {
		path, err := os.Getwd()
		utils.PanicOnError(err)
		path, err = filepath.Abs(path)
		utils.PanicOnError(err)
		configPath = filepath.Join(path, "src/config/default.toml")
	}

	config = &Config{
		InfoDir:              "~/.skizze",
		DataDir:              "~/.skizze/data",
		Port:                 3956,
		SaveThresholdSeconds: 1,
	}
	_, err := os.Open(configPath)
	if err != nil {
		logger.Warningf("Unable to find config file, using defaults")
		return config
	}
	if _, err := toml.DecodeFile(configPath, &config); err != nil {
		utils.PanicOnError(err)
	}
	return config
}

// GetConfig returns a singleton Configuration
func GetConfig() *Config {
	if config == nil {
		config = parseConfigTOML()
		usr, err := user.Current()
		utils.PanicOnError(err)
		dir := usr.HomeDir

		infoDir := strings.TrimSpace(os.Getenv("SKZ_INFO_DIR"))
		if len(infoDir) == 0 {
			if config.InfoDir[:2] == "~/" {
				infoDir = strings.Replace(config.InfoDir, "~", dir, 1)
			}
		}

		dataDir := strings.TrimSpace(os.Getenv("SKZ_DATA_DIR"))
		if len(dataDir) == 0 {
			if config.DataDir[:2] == "~/" {
				dataDir = strings.Replace(config.DataDir, "~", dir, 1)
			}
		}

		portInt, err := strconv.Atoi(strings.TrimSpace(os.Getenv("SKZ_PORT")))
		port := uint(portInt)
		if err != nil {
			port = config.Port
		}

		saveThresholdSecondsInt, err := strconv.Atoi(strings.TrimSpace(os.Getenv("SKZ_SAVE_TRESHOLD_SECS")))
		saveThresholdSeconds := uint(saveThresholdSecondsInt)
		if err != nil {
			saveThresholdSeconds = config.SaveThresholdSeconds
		}

		if saveThresholdSeconds < 3 {
			saveThresholdSeconds = 3
		}

		if err := os.MkdirAll(dataDir, os.ModePerm); err != nil {
			panic(err)
		}

		config = &Config{
			infoDir,
			dataDir,
			port,
			saveThresholdSeconds,
		}
	}
	return config
}

// Reset ...
func Reset() {
	config = nil
}
