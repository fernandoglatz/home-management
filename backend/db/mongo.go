package db

import (
	"context"
	"errors"
	"log"
	"strings"

	"github.com/fernandoglatz/home-management/models"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mongodb"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client
var userCollection *mongo.Collection

// ConnectToMongoDB establishes a connection to MongoDB
func ConnectToMongoDB(uri string, database string) error {
	// Set up MongoDB connection

	clientOptions := options.Client().ApplyURI(uri)
	var err error
	client, err = mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return err
	}

	// Check the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		return err
	}

	log.Println("Connected to MongoDB!")

	// Run migrations
	driver, err := mongodb.WithInstance(client, &mongodb.Config{
		DatabaseName: database,
	})
	if err != nil {
		return err
	}

	migrations, err := migrate.NewWithDatabaseInstance("file://resources/db/migrations", database, driver)
	if err != nil {
		return err
	}

	if err := migrations.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	// Set the user collection
	userCollection = client.Database(database).Collection("users")

	return nil
}

// InsertUser inserts a new user into MongoDB
func InsertUser(ctx context.Context, user *models.User) error {
	uuidObj, err := uuid.NewRandom()
	uuidStr := uuidObj.String()
	user.ID = strings.Replace(uuidStr, "-", "", -1)

	_, err = userCollection.InsertOne(ctx, user)
	return err
}

// GetUser retrieves a user from MongoDB by ID
func GetUser(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	err := userCollection.FindOne(ctx, bson.M{"id": id}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// getAllUsers retrieves all users from MongoDB
func GetAllUsers(ctx context.Context) ([]models.User, error) {
	cursor, err := userCollection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []models.User = []models.User{}
	for cursor.Next(ctx) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

// UpdateUser updates a user in MongoDB
func UpdateUser(ctx context.Context, user *models.User) error {
	filter := bson.M{"id": user.ID}
	_, err := userCollection.ReplaceOne(ctx, filter, user)
	return err
}

// DeleteUser deletes a user from MongoDB by ID
func DeleteUser(ctx context.Context, id string) error {
	result, err := userCollection.DeleteOne(ctx, bson.M{"id": id})

	if result.DeletedCount == 0 {
		return errors.New("User not found")
	}

	return err
}

// FindUserByID retrieves a user from MongoDB by ID
func FindUserByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User

	err := userCollection.FindOne(ctx, bson.M{"id": id}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
