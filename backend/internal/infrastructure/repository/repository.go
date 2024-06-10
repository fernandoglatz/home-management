package repository

import (
	"context"
	"strings"
	"sync"
	"time"

	"fernandoglatz/home-management/internal/core/common/utils"
	"fernandoglatz/home-management/internal/core/common/utils/constants"
	"fernandoglatz/home-management/internal/core/common/utils/exceptions"
	"fernandoglatz/home-management/internal/core/entity"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var repositories map[string]any
var repositoryMutex sync.Mutex

type Repository[T entity.IEntity] struct {
	collection mongo.Collection
}

func GetGenericRepository[T entity.IEntity]() Repository[T] {
	entity := utils.Instance[T]()
	typeName := utils.GetTypeName(entity)

	repositoryMutex.Lock()
	defer repositoryMutex.Unlock()

	if repositories == nil {
		repositories = make(map[string]any)
	}

	repository := repositories[typeName]

	if repository == nil {
		collectionName := entity.GetCollectionName()
		repository = Repository[T]{
			collection: utils.MongoDatabase.GetCollection(collectionName),
		}

		repositories[typeName] = repository
	}

	return repository.(Repository[T])
}

func (repository Repository[T]) Get(ctx context.Context, id string) (T, *exceptions.WrappedError) {
	filter := bson.M{"id": id}
	return repository.getByFilter(ctx, filter)
}

func (repository Repository[T]) getByFilter(ctx context.Context, filter any) (T, *exceptions.WrappedError) {
	entity := utils.Instance[T]()

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

	repository.CorrecTimezone(entity)
	return entity, nil
}

func (repository Repository[T]) GetAll(ctx context.Context) ([]T, *exceptions.WrappedError) {
	var entities []T = []T{}

	cursor, err := repository.collection.Find(ctx, bson.M{})
	if err != nil {
		return entities, &exceptions.WrappedError{
			Error: err,
		}
	}

	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		entity := utils.Instance[T]()
		err = cursor.Decode(&entity)
		if err != nil {
			return entities, &exceptions.WrappedError{
				Error: err,
			}
		}

		repository.CorrecTimezone(entity)
		entities = append(entities, entity)
	}

	return entities, nil
}

func (repository Repository[T]) Save(ctx context.Context, entity T) *exceptions.WrappedError {
	now := time.Now()
	entity.SetUpdatedAt(now)

	if len(entity.GetID()) == constants.ZERO {
		uuidObj, _ := uuid.NewRandom()
		uuidStr := uuidObj.String()
		entity.SetID(strings.Replace(uuidStr, "-", constants.EMPTY, -1))
	}

	var err error

	if entity.GetCreatedAt().IsZero() {
		entity.SetCreatedAt(now)
		_, err = repository.collection.InsertOne(ctx, entity)
	} else {
		filter := bson.M{"id": entity.GetID()}
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

func (repository Repository[T]) Remove(ctx context.Context, entity T) *exceptions.WrappedError {
	filter := bson.M{"id": entity.GetID()}
	_, err := repository.collection.DeleteOne(ctx, filter)
	if err != nil {
		return &exceptions.WrappedError{
			Error: err,
		}
	}

	return nil
}

func (repository Repository[T]) CorrecTimezone(entity T) {
	location, _ := time.LoadLocation(utils.GetTimezone())

	createdAt := entity.GetCreatedAt()
	updatedAt := entity.GetUpdatedAt()

	entity.SetCreatedAt(createdAt.In(location))
	entity.SetUpdatedAt(updatedAt.In(location))
}
