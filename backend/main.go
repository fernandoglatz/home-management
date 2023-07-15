package main

import (
	"log"

	"github.com/fernandoglatz/home-management/api"
	"github.com/fernandoglatz/home-management/configs"
	"github.com/fernandoglatz/home-management/utils"
)

func main() {
	// Load configuration
	err := configs.LoadConfig()
	if err != nil {
		log.Fatalf(err.Error())
	}

	// Connect to MongoDB
	err = utils.ConnectToMongoDB()
	if err != nil {
		log.Fatalf(err.Error())
	}

	err = api.Setup()
	if err != nil {
		log.Fatalf(err.Error())
	}

}
