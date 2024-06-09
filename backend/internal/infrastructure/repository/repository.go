package repository

import (
	"context"
	"strings"
	"time"

	"fernandoglatz/home-management/internal/core/common/utils"
	"fernandoglatz/home-management/internal/core/common/utils/constants"
	"fernandoglatz/home-management/internal/core/common/utils/exceptions"
	"fernandoglatz/home-management/internal/core/entity"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository[T entity.IEntity] struct {
	collection *mongo.Collection
}

func NewRepository[T entity.IEntity](baseEntity T) *Repository[T] {
	collectionName := baseEntity.GetCollectionName()
	return &Repository[T]{
		collection: utils.MongoDatabase.GetCollection(collectionName),
	}
}

func (repository *Repository[T]) Get(ctx context.Context, id string) (T, *exceptions.WrappedError) {
	filter := bson.M{"id": id}
	return repository.getByFilter(ctx, filter)
}

func (repository *Repository[T]) getByFilter(ctx context.Context, filter any) (T, *exceptions.WrappedError) {
	var entity T

	err := repository.collection.FindOne(ctx, filter).Decode(&entity)
	if err == mongo.ErrNoDocuments {
		return entity, &exceptions.WrappedError{
			BaseError: exceptions.RecordNotFound,
		}
	} else if err != nil {
		return entity, &exceptions.WrappedError{
			Error: err,
		}
	}

	repository.correctTimezone(&entity)
	return entity, nil
}

func (repository *Repository[T]) GetAll(ctx context.Context) ([]T, *exceptions.WrappedError) {
	var entitys []T = []T{}

	cursor, err := repository.collection.Find(ctx, bson.D{})
	if err != nil {
		return entitys, &exceptions.WrappedError{
			Error: err,
		}
	}

	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var entity T
		err = cursor.Decode(&entity)
		if err != nil {
			return entitys, &exceptions.WrappedError{
				Error: err,
			}
		}

		repository.correctTimezone(&entity)
		entitys = append(entitys, entity)
	}

	return entitys, nil
}

func (repository *Repository[T]) Save(ctx context.Context, entity *T) *exceptions.WrappedError {
	now := time.Now()
	(*entity).SetUpdatedAt(now)

	if len((*entity).GetID()) == constants.ZERO {
		uuidObj, _ := uuid.NewRandom()
		uuidStr := uuidObj.String()
		(*entity).SetID(strings.Replace(uuidStr, "-", "", -1))
	}

	var err error

	if (*entity).GetCreatedAt().IsZero() {
		(*entity).SetCreatedAt(now)
		_, err = repository.collection.InsertOne(ctx, entity)
	} else {
		filter := bson.M{"id": (*entity).GetID()}
		_, err = repository.collection.ReplaceOne(ctx, filter, entity)
	}

	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return &exceptions.WrappedError{
				BaseError: exceptions.DuplicatedRecord,
			}

		} else {
			return &exceptions.WrappedError{
				Error: err,
			}
		}
	}

	return nil
}

func (repository *Repository[T]) Remove(ctx context.Context, entity T) *exceptions.WrappedError {
	filter := bson.M{"id": entity.GetID()}
	_, err := repository.collection.DeleteOne(ctx, filter)
	if err != nil {
		return &exceptions.WrappedError{
			Error: err,
		}
	}

	return nil
}
func (repository *Repository[T]) correctTimezone(entity *T) {
	location, _ := time.LoadLocation(utils.GetTimezone())

	createdAt := (*entity).GetCreatedAt()
	updatedAt := (*entity).GetUpdatedAt()

	(*entity).SetCreatedAt(createdAt.In(location))
	(*entity).SetUpdatedAt(updatedAt.In(location))
}
