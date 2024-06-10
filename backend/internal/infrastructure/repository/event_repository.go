package repository

import (
	"context"
	"fernandoglatz/home-management/internal/core/common/utils"
	"fernandoglatz/home-management/internal/core/common/utils/exceptions"
	"fernandoglatz/home-management/internal/core/entity"
	"sync"
)

var eventRepositories map[string]any
var eventRepositoryMutex sync.Mutex

type EventRepository[T entity.IEntity] struct {
	repository Repository[T]
}

func GetEventRepository[T entity.IEntity]() EventRepository[T] {
	var entity T
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

func (eventRepository *EventRepository[T]) Get(ctx context.Context, id string) (T, *exceptions.WrappedError) {
	return eventRepository.repository.Get(ctx, id)
}

func (eventRepository *EventRepository[T]) GetAll(ctx context.Context) ([]T, *exceptions.WrappedError) {
	return eventRepository.repository.GetAll(ctx)
}

func (eventRepository *EventRepository[T]) Save(ctx context.Context, entity *T) *exceptions.WrappedError {
	return eventRepository.repository.Save(ctx, entity)
}
func (eventRepository *EventRepository[T]) Remove(ctx context.Context, entity T) *exceptions.WrappedError {
	return eventRepository.repository.Remove(ctx, entity)
}

func (eventRepository *EventRepository[T]) CorrecTimezone(entity *T) {
	eventRepository.repository.CorrecTimezone(entity)
}
