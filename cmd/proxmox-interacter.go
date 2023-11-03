package main

import (
	"main/pkg"
	"main/pkg/app"
	"main/pkg/logger"

	"github.com/spf13/cobra"
)

var version = "unknown"

func Execute(configPath string) {
	if configPath == "" {
		logger.GetDefaultLogger().Fatal().Msg("Cannot start without config")
	}

	config := pkg.LoadConfig(configPath)
	newApp := app.NewApp(config, version)
	newApp.Start()
}

func main() {
	var configPath string

	rootCmd := &cobra.Command{
		Use:     "proxmox-interacter",
		Long:    "A Telegram bot to interact with your Proxmox instances.",
		Version: version,
		Run: func(cmd *cobra.Command, args []string) {
			Execute(configPath)
		},
	}

	rootCmd.PersistentFlags().StringVar(&configPath, "config", "", "Config file path")
	if err := rootCmd.MarkPersistentFlagRequired("config"); err != nil {
		logger.GetDefaultLogger().Fatal().Err(err).Msg("Could not set flag as required")
	}

	if err := rootCmd.Execute(); err != nil {
		logger.GetDefaultLogger().Fatal().Err(err).Msg("Could not start application")
	}
}
