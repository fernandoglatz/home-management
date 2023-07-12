package main

import (
	"io/ioutil"
	"log"

	"github.com/fernandoglatz/home-management/db"
	"github.com/fernandoglatz/home-management/handlers"
	"github.com/fernandoglatz/home-management/models"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
)

var config models.Config

func main() {
	// Load configuration
	err := loadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Connect to MongoDB
	err = db.ConnectToMongoDB(config.Database.URI, config.Database.Name)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Initialize Gin router
	router := gin.Default()

	// Define API routes
	router.POST("/users", handlers.CreateUser)
	router.GET("/users/:id", handlers.GetUser)
	router.GET("/users", handlers.GetAllUsers)
	router.PUT("/users/:id", handlers.UpdateUser)
	router.DELETE("/users/:id", handlers.DeleteUser)

	// Start the server
	log.Fatal(router.Run(config.Server.Listening))
}

func loadConfig() error {
	data, err := ioutil.ReadFile("resources/config/config.yml")
	if err != nil {
		log.Fatalf("Failed to read configuration file: %v", err)
	}

	// Parse the YAML data into Config struct
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("Failed to parse configuration file: %v", err)
	}

	// Access the configuration values
	log.Println("Database URI:", config.Database.URI)
	log.Println("Database Name:", config.Database.Name)
	log.Println("Server Listening:", config.Server.Listening)

	return nil
}
