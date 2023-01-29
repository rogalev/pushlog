package main

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rogalev/pushlog/internal/config"
	"github.com/rogalev/pushlog/internal/logging"
	"github.com/rogalev/pushlog/internal/publisher"
	"github.com/rogalev/pushlog/internal/receiver"
	"github.com/rogalev/pushlog/internal/service"
	"github.com/rogalev/pushlog/internal/storage"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"os"
)

func main() {

	var configFile string

	var rootCmd = &cobra.Command{
		Use:   "",
		Short: "PushLog",
		Long:  "PushLog - service for collecting and sending logs",
		Run: func(cmd *cobra.Command, args []string) {
			cfg, err := config.Init(configFile)

			if err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}

			logging.SetupConfig(cfg)
			logger := logging.GetInstance()

			r, err := receiver.NewInstance(cfg)
			if err != nil {
				logger.Error("wrong receiver instance", zap.Error(err))
			}

			p, err := publisher.NewInstance(cfg)
			if err != nil {
				logger.Error("wrong publisher instance", zap.Error(err))
			}

			s, err := storage.NewInstance(cfg)
			if err != nil {
				logger.Error("wrong storage instance", zap.Error(err))
			}

			srv := service.NewInstance(r, p, s)
			srv.Run()
		},
	}

	rootCmd.Flags().StringVarP(&configFile, "config", "c", "config.json", "path to config file")

	var outputConfigFile string

	var genConfigCmd = &cobra.Command{
		Use:   "genconfig",
		Short: "Generate config file with default values",
		Long:  "Generate config file with default values",
		Run: func(cmd *cobra.Command, args []string) {

			if err := config.GenerateFileWithDefaultValues(outputConfigFile); err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}

			fmt.Printf("Config file successfully generated!\n")
			os.Exit(1)
		},
	}

	genConfigCmd.Flags().StringVarP(&outputConfigFile, "config", "c", "config.json", "path to output config file")

	var telegramToken string
	var telegramLastUpdate int

	var tgUpdatesCmd = &cobra.Command{
		Use:   "tgupdates",
		Short: "Watch telegram bot updates",
		Long:  "Watch telegram bot updates",
		Run: func(cmd *cobra.Command, args []string) {

			bot, err := tgbotapi.NewBotAPI(telegramToken)

			if err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}

			updateConfig := tgbotapi.NewUpdate(telegramLastUpdate)
			updateConfig.Timeout = 30
			updates := bot.GetUpdatesChan(updateConfig)

			for update := range updates {
				data, err := json.Marshal(update)

				if err != nil {
					fmt.Printf("Error: %v\n", err)
					os.Exit(1)
				}

				fmt.Printf("%s\n", data)
			}
		},
	}

	tgUpdatesCmd.Flags().StringVarP(&telegramToken, "token", "t", "", "Telegram bot token")
	tgUpdatesCmd.Flags().IntVarP(&telegramLastUpdate, "update", "u", 0, "Telegram update offset")

	rootCmd.AddCommand(genConfigCmd)
	rootCmd.AddCommand(tgUpdatesCmd)

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}

}
