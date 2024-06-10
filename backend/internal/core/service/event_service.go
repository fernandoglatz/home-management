package service

import (
	"context"
	"fernandoglatz/home-management/internal/core/common/utils/exceptions"
	"fernandoglatz/home-management/internal/core/entity"
	"fernandoglatz/home-management/internal/infrastructure/repository"
	"sync"
)

var eventService any
var eventServiceMutex sync.Mutex

type EventService[T entity.IEntity] struct {
	service    Service[T]
	repository repository.EventRepository[T]
}

func GetEventService[T entity.IEntity]() EventService[T] {
	eventServiceMutex.Lock()
	defer eventServiceMutex.Unlock()

	if eventService == nil {
		service := GetGenericService[T]()
		repository := repository.GetEventRepository[T]()

		eventService = EventService[T]{
			service:    service,
			repository: repository,
		}
	}

	return eventService.(EventService[T])
}

func (eventService *EventService[T]) Get(ctx context.Context, id string) (T, *exceptions.WrappedError) {
	return eventService.service.Get(ctx, id)
}

func (eventService *EventService[T]) GetAll(ctx context.Context) ([]T, *exceptions.WrappedError) {
	return eventService.service.GetAll(ctx)
}

func (eventService *EventService[T]) Save(ctx context.Context, entity *T) *exceptions.WrappedError {
	return eventService.service.Save(ctx, entity)
}

func (eventService *EventService[T]) Remove(ctx context.Context, entity T) *exceptions.WrappedError {
	return eventService.service.Remove(ctx, entity)
}
