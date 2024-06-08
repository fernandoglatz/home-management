package service

import (
	"context"

	"fernandoglatz/home-management/internal/core/common/utils/exceptions"
	"fernandoglatz/home-management/internal/core/entity"
	"fernandoglatz/home-management/internal/infrastructure/repository"
)

type Service[T entity.IEntity] struct {
	repository repository.Repository[T]
}

func NewService[T entity.IEntity](repository repository.Repository[T]) Service[T] {
	return Service[T]{
		repository: repository,
	}
}

func (service *Service[T]) Get(ctx context.Context, id string) (T, *exceptions.WrappedError) {
	return service.repository.Get(ctx, id)
}

func (service *Service[T]) GetAll(ctx context.Context) ([]T, *exceptions.WrappedError) {
	return service.repository.GetAll(ctx)
}

func (service *Service[T]) Save(ctx context.Context, entity *T) *exceptions.WrappedError {
	return service.repository.Save(ctx, entity)
}

func (service *Service[T]) Remove(ctx context.Context, entity T) *exceptions.WrappedError {
	return service.repository.Remove(ctx, entity)
}
