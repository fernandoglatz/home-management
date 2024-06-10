package service

import (
	"context"
	"fernandoglatz/home-management/internal/core/common/utils/exceptions"
	"fernandoglatz/home-management/internal/core/entity"
)

type IService[T entity.IEntity] interface {
	Get(ctx context.Context, id string) (T, *exceptions.WrappedError)
	GetAll(ctx context.Context) ([]T, *exceptions.WrappedError)
	Save(ctx context.Context, entity *T) *exceptions.WrappedError
	Remove(ctx context.Context, entity T) *exceptions.WrappedError
}
