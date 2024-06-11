package service

import (
	"context"
	"sync"

	"fernandoglatz/home-management/internal/core/common/utils"
	"fernandoglatz/home-management/internal/core/common/utils/exceptions"
	"fernandoglatz/home-management/internal/core/entity"
	repository_port "fernandoglatz/home-management/internal/core/port/repository"
	"fernandoglatz/home-management/internal/infrastructure/repository"
)

var services map[string]any
var serviceMutex sync.Mutex

type Service[T entity.IEntity] struct {
	repository repository_port.IRepository[T]
}

func GetGenericService[T entity.IEntity]() Service[T] {
	repository := repository.GetGenericRepository[T]()
	return GetService[T](repository)
}

func GetService[T entity.IEntity](repository repository_port.IRepository[T]) Service[T] {
	var entity T
	typeName := utils.GetTypeName(entity)

	serviceMutex.Lock()
	defer serviceMutex.Unlock()

	if services == nil {
		services = make(map[string]any)
	}

	service := services[typeName]

	if service == nil {
		service = Service[T]{
			repository: repository,
		}

		services[typeName] = service
	}

	return service.(Service[T])
}

func (service Service[T]) Get(ctx context.Context, id string) (T, *exceptions.WrappedError) {
	return service.repository.Get(ctx, id)
}

func (service Service[T]) GetAll(ctx context.Context, page int, limit int) ([]T, *exceptions.WrappedError) {
	return service.repository.GetAll(ctx, page, limit)
}

func (service Service[T]) Save(ctx context.Context, entity T) *exceptions.WrappedError {
	return service.repository.Save(ctx, entity)
}

func (service Service[T]) Remove(ctx context.Context, entity T) *exceptions.WrappedError {
	return service.repository.Remove(ctx, entity)
}
