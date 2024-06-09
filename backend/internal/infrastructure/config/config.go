package config

import (
	"context"
	"errors"
	"fernandoglatz/home-management/internal/core/common/utils/constants"
	"fernandoglatz/home-management/internal/core/common/utils/log"
	"fernandoglatz/home-management/internal/infrastructure/config/format"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Listening   string `yaml:"listening"`
		ContextPath string `yaml:"context-path"`
	} `yaml:"server"`

	Data struct {
		Mongo struct {
			Uri      string `yaml:"uri"`
			Database string `yaml:"database"`
		} `yaml:"mongo"`

		Redis struct {
			Address  string `yaml:"address"`
			Password string `yaml:"password"`
			Db       int    `yaml:"db"`

			TTL struct {
				RfEvent time.Duration `yaml:"rf-event"`
			} `yaml:"ttl"`
		} `yaml:"redis"`
	} `yaml:"data"`

	Broker struct {
		Mqtt struct {
			Uri      string `yaml:"uri"`
			ClientId string `yaml:"client-id"`
			User     string `yaml:"user"`
			Password string `yaml:"password"`

			Topics struct {
				Broadcast string `yaml:"broadcast"`
				Devices   string `yaml:"devices"`
			} `yaml:"topics"`
		} `yaml:"mqtt"`

		RabbitMQ struct {
			Uri      string `yaml:"uri"`
			User     string `yaml:"user"`
			Password string `yaml:"password"`

			Exchanges struct {
				Events string `yaml:"events"`
			} `yaml:"exchanges"`

			Queues struct {
				Events string `yaml:"events"`
			} `yaml:"queues"`
		} `yaml:"rabbitmq"`
	} `yaml:"broker"`

	Log struct {
		Level   string        `yaml:"level"`
		Format  format.Format `yaml:"format"`
		Colored bool          `yaml:"colored"`
	} `yaml:"log"`
}

var ApplicationConfig Config

func LoadConfig(ctx context.Context) error {
	loadProfile(ctx)

	err := loadLocalConfig(ctx)
	if err != nil {
		return err
	}

	logConfig := ApplicationConfig.Log
	log.ReconfigureLogger(ctx, logConfig.Format, logConfig.Level, logConfig.Colored)

	return nil
}

func IsDevProfile() bool {
	profile := os.Getenv(constants.PROFILE)
	return constants.DEV_PROFILE == profile
}

func loadProfile(ctx context.Context) {
	profile := os.Getenv(constants.PROFILE)
	if len(profile) == constants.ZERO {
		profile = constants.DEV_PROFILE
		os.Setenv(constants.PROFILE, profile)
	}

	log.SetupLogger(profile) //after setup profile
	log.Info(ctx).Msg("Profile loaded: " + profile)
}

func loadLocalConfig(ctx context.Context) error {
	log.Info(ctx).Msg("Loading local config")

	data, err := os.ReadFile("conf/application.yml")
	if err != nil {
		return errors.New("Failed to read configuration file: " + err.Error())
	}

	err = yaml.Unmarshal(data, &ApplicationConfig)
	if err != nil {
		return errors.New("Failed to parse configuration file: " + err.Error())
	}

	log.Info(ctx).Msg("Loaded local config")

	return nil
}
