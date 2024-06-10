package repository

import (
	"context"
	"fernandoglatz/home-management/internal/core/common/utils/exceptions"
	"fernandoglatz/home-management/internal/core/entity"
	"sync"
)

var eventRepository any
var eventRepositoryMutex sync.Mutex

type EventRepository[T entity.IEntity] struct {
	repository Repository[T]
}

func GetEventRepository[T entity.IEntity]() EventRepository[T] {
	eventRepositoryMutex.Lock()
	defer eventRepositoryMutex.Unlock()

	if eventRepository == nil {
		repository := GetGenericRepository[T]()

		eventRepository = EventRepository[T]{
			repository: repository,
		}
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
