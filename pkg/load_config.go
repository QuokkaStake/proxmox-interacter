package pkg

import (
	"main/pkg/config"
	"main/pkg/logger"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/creasty/defaults"
)

func LoadConfig(path string) *config.Config {
	configString, err := os.ReadFile(path)
	if err != nil {
		logger.GetDefaultLogger().Fatal().Err(err).Msg("Could not read config file")
	}

	var configStruct *config.Config

	if _, err = toml.Decode(string(configString), &configStruct); err != nil {
		logger.GetDefaultLogger().Fatal().Err(err).Msg("Could not unmarshal config file")
	}

	if err = defaults.Set(configStruct); err != nil {
		logger.GetDefaultLogger().Fatal().Err(err).Msg("Could not set default config values")
	}

	return configStruct
}
