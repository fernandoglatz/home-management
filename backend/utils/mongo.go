package utils

import (
	"context"
	"log"

	"github.com/fernandoglatz/home-management/backend/configs"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mongodb"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var database *mongo.Database
var collections map[string]*mongo.Collection = make(map[string]*mongo.Collection)

func ConnectToMongoDB() error {
	config := configs.ApplicationConfig
	uri := config.Database.URI
	databaseName := config.Database.Name

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return err
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		return err
	}

	log.Println("Connected to MongoDB!")

	mongodbDriver, err := mongodb.WithInstance(client, &mongodb.Config{
		DatabaseName: databaseName,
	})
	if err != nil {
		return err
	}

	migrations, err := migrate.NewWithDatabaseInstance("file://scripts/mongo/migrations", databaseName, mongodbDriver)
	if err != nil {
		return err
	}

	err = migrations.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	database = client.Database(databaseName)
	return nil
}

func GetMongoDbCollection(collectionName string) *mongo.Collection {
	collection := collections[collectionName]

	if collection == nil {
		collection = database.Collection(collectionName)
		collections[collectionName] = collection
	}

	return collection
}
