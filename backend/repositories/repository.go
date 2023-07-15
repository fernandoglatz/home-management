package repositories

import (
	"context"
	"strings"

	"github.com/fernandoglatz/home-management/models"
	"github.com/fernandoglatz/home-management/utils"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository[T models.IEntity] struct {
	baseEntity T
}

func NewRepository[T models.IEntity]() *Repository[T] {
	return &Repository[T]{}
}

func (repository *Repository[T]) Insert(ctx context.Context, entity T) error {
	uuidObj, err := uuid.NewRandom()
	uuidStr := uuidObj.String()
	id := strings.Replace(uuidStr, "-", "", -1)
	entity.SetID(id)

	collection := GetCollection(repository.baseEntity)
	_, err = collection.InsertOne(ctx, entity)
	return err
}

func (repository *Repository[T]) Update(ctx context.Context, entity models.IEntity) error {
	filter := bson.M{"id": entity.GetID}

	collection := GetCollection(repository.baseEntity)
	_, err := collection.ReplaceOne(ctx, filter, entity)
	return err
}

func (repository *Repository[T]) Delete(ctx context.Context, id string) error {
	collection := GetCollection(repository.baseEntity)
	_, err := collection.DeleteOne(ctx, bson.M{"id": id})
	return err
}

func (repository *Repository[T]) FindByID(ctx context.Context, id string) (T, error) {
	var value T
	collection := GetCollection(repository.baseEntity)

	err := collection.FindOne(ctx, bson.M{"id": id}).Decode(&value)
	if err != nil {
		return value, err
	}
	return value, nil
}

func (repository *Repository[T]) FindAll(ctx context.Context) ([]T, error) {
	collection := GetCollection(repository.baseEntity)

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
