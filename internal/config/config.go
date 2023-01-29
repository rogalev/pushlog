package config

import (
	"encoding/json"
	"errors"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"path/filepath"
)

type (
	HttpReceiverConfig struct {
		Host         string `json:"host"`
		Port         int    `json:"port"`
		ReadTimeout  int    `json:"readTimeout"`
		WriteTimeout int    `json:"WriteTimeout"`
	}

	TelegramPublisherConfig struct {
		Token      string `json:"token"`
		Chat       int64  `json:"chat"`
		SendAsFile bool   `json:"sendAsFile"`
		TmpFileDir string `json:"tmpFileDir"`
	}

	MemoryStorageConfig struct {
		ExpirationGCPeriod int64 `json:"expirationGCPeriod"`
	}

	Config struct {
		ReceiverEngine  string `json:"receiverEngine"`
		PublisherEngine string `json:"publisherEngine"`
		StorageEngine   string `json:"storageEngine"`

		LogFile string `json:"logFile"`

		HttpReceiver      HttpReceiverConfig      `json:"httpReceiver"`
		TelegramPublisher TelegramPublisherConfig `json:"telegramPublisher"`
		MemoryStorage     MemoryStorageConfig     `json:"memoryStorage"`
	}
)

func getDefaultConfig() Config {
	return Config{
		ReceiverEngine:  "http",
		PublisherEngine: "telegram",
		StorageEngine:   "memory",

		LogFile: "./debug.log",

		HttpReceiver: HttpReceiverConfig{
			Host:         "localhost",
			Port:         8000,
			ReadTimeout:  10,
			WriteTimeout: 10,
		},

		TelegramPublisher: TelegramPublisherConfig{
			Token:      "",
			Chat:       0,
			SendAsFile: true,
			TmpFileDir: "/tmp",
		},

		MemoryStorage: MemoryStorageConfig{
			ExpirationGCPeriod: 10,
		},
	}
}

func Init(configFile string) (Config, error) {

	cfg := getDefaultConfig()

	ext := filepath.Ext(configFile)

	if ext != ".json" {
		return cfg, errors.New("invalid config file extension")
	}

	viper.SetConfigFile(configFile)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return cfg, err
		}
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}

func GenerateFileWithDefaultValues(outputConfigFile string) error {

	var cfg = getDefaultConfig()

	_, err := os.Stat(outputConfigFile)

	if err == nil {
		return errors.New("config file already exists")
	}
	if !os.IsNotExist(err) {
		return err
	}

	ext := filepath.Ext(outputConfigFile)

	if ext != ".json" {
		return errors.New("invalid config file extension")
	}

	file, err := json.MarshalIndent(cfg, "", "    ")

	if err != nil {
		return err
	}

	if err = ioutil.WriteFile(outputConfigFile, file, 0644); err != nil {
		return err
	}

	return nil
}
