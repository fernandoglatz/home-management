package repository

import (
	"context"
	"fernandoglatz/home-management/internal/core/common/utils"
	"fernandoglatz/home-management/internal/core/common/utils/exceptions"
	"fernandoglatz/home-management/internal/core/entity"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

var eventRepositories map[string]any
var eventRepositoryMutex sync.Mutex

type EventRepository[T entity.IEvent] struct {
	repository Repository[T]
}

func GetEventRepository[T entity.IEvent]() EventRepository[T] {
	entity := utils.Instance[T]()
	typeName := utils.GetTypeName(entity)

	eventRepositoryMutex.Lock()
	defer eventRepositoryMutex.Unlock()

	if eventRepositories == nil {
		eventRepositories = make(map[string]any)
	}

	eventRepository := eventRepositories[typeName]

	if eventRepository == nil {
		repository := GetGenericRepository[T]()

		eventRepository = EventRepository[T]{
			repository: repository,
		}

		eventRepositories[typeName] = eventRepository
	}

	return eventRepository.(EventRepository[T])
}

func (eventRepository EventRepository[T]) Get(ctx context.Context, id string) (T, *exceptions.WrappedError) {
	entity, err := eventRepository.repository.Get(ctx, id)

	if err == nil {
		eventRepository.CorrecTimezone(entity)
	}

	return entity, err
}

func (eventRepository EventRepository[T]) GetAll(ctx context.Context) ([]T, *exceptions.WrappedError) {
	entities, err := eventRepository.repository.GetAll(ctx)

	for _, entity := range entities {
		eventRepository.CorrecTimezone(entity)
	}

	return entities, err
}

func (eventRepository EventRepository[T]) Save(ctx context.Context, entity T) *exceptions.WrappedError {
	return eventRepository.repository.Save(ctx, entity)
}
func (eventRepository EventRepository[T]) Remove(ctx context.Context, entity T) *exceptions.WrappedError {
	return eventRepository.repository.Remove(ctx, entity)
}

func (eventRepository EventRepository[T]) CorrecTimezone(entity T) {
	location, _ := time.LoadLocation(utils.GetTimezone())

	date := entity.GetDate()
	entity.SetDate(date.In(location))

	eventRepository.repository.CorrecTimezone(entity)
}

func (eventRepository EventRepository[T]) GetRfEvents(ctx context.Context, code int, bits int, protocol int, frequency int, startDate time.Time, endDate time.Time) ([]entity.RfEvent, *exceptions.WrappedError) {
	var events []entity.RfEvent = []entity.RfEvent{}

	rfEventRepository := GetEventRepository[*entity.RfEvent]()
	filter := bson.M{
		"code":      code,
		"bits":      bits,
		"protocol":  protocol,
		"frequency": frequency,
		"date": bson.M{
			"$gte": startDate,
			"$lte": endDate,
		},
	}

	cursor, err := rfEventRepository.repository.collection.Find(ctx, filter)
	if err != nil {
		return events, &exceptions.WrappedError{
			Error: err,
		}
	}

	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var event entity.RfEvent
		err = cursor.Decode(&event)
		if err != nil {
			return events, &exceptions.WrappedError{
				Error: err,
			}
		}

		rfEventRepository.repository.CorrecTimezone(&event)
		events = append(events, event)
	}

	return events, nil
}
