package repositories

import (
	"context"
	"strings"

	"github.com/fernandoglatz/home-management/backend/models"
	"github.com/fernandoglatz/home-management/backend/utils"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository[T models.IEntity] struct {
	baseEntity T
}

func NewRepository[T models.IEntity]() Repository[T] {
	return Repository[T]{}
}

func (repository *Repository[T]) Insert(entity T) error {
	uuidObj, err := uuid.NewRandom()
	uuidStr := uuidObj.String()
	id := strings.Replace(uuidStr, "-", "", -1)
	entity.SetID(id)

	collection := GetCollection(repository.baseEntity)
	_, err = collection.InsertOne(context.Background(), entity)
	return err
}

func (repository *Repository[T]) Update(entity T) error {
	id := entity.GetID()
	filter := bson.M{"id": id}
	collection := GetCollection(repository.baseEntity)

	_, err := collection.ReplaceOne(context.Background(), filter, entity)
	return err
}

func (repository *Repository[T]) Delete(entity T) error {
	id := entity.GetID()
	filter := bson.M{"id": id}
	collection := GetCollection(repository.baseEntity)

	_, err := collection.DeleteOne(context.Background(), filter)
	return err
}

func (repository *Repository[T]) FindByID(id string) (T, error) {
	var value T

	filter := bson.M{"id": id}
	collection := GetCollection(repository.baseEntity)

	err := collection.FindOne(context.Background(), filter).Decode(&value)
	if err != nil {
		return value, err
	}
	return value, nil
}

func (repository *Repository[T]) FindAll() ([]T, error) {
	collection := GetCollection(repository.baseEntity)

	ctx := context.Background()
	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var values []T = []T{}
	for cursor.Next(ctx) {
		var value T
		err = cursor.Decode(&value)
		if err != nil {
			return nil, err
		}
		values = append(values, value)
	}

	return values, nil
}

func GetCollection(ientity models.IEntity) *mongo.Collection {
	return utils.GetMongoDbCollection(ientity.GetCollectionName())
}
