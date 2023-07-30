package configs

import (
	"errors"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Database struct {
		URI  string `yaml:"uri"`
		Name string `yaml:"name"`
	} `yaml:"database"`
	Server struct {
		Listening   string `yaml:"listening"`
		AllowOrigin string `yaml:"allow-origin"`
	} `yaml:"server"`
	Mosquitto struct {
		ClientID   string `yaml:"clientID"`
		URI   string `yaml:"uri"`
		Topic string `yaml:"topic"`
	} `yaml:"mosquitto"`
}

var ApplicationConfig Config

func LoadConfig() error {
	data, err := ioutil.ReadFile("configs/config.yml")
	if err != nil {
		return errors.New("Failed to read configuration file: " + err.Error())
	}

	err = yaml.Unmarshal(data, &ApplicationConfig)
	if err != nil {
		return errors.New("Failed to parse configuration file: " + err.Error())
	}

	log.Println("Database URI:", ApplicationConfig.Database.URI)
	log.Println("Database Name:", ApplicationConfig.Database.Name)
	log.Println("Server Listening:", ApplicationConfig.Server.Listening)

	return nil
}
