package repository

import (
	"context"
	"fernandoglatz/home-management/internal/core/common/utils/exceptions"
	"fernandoglatz/home-management/internal/core/entity"
)

type IRepository[T entity.IEntity] interface {
	Get(ctx context.Context, id string) (T, *exceptions.WrappedError)
	GetAll(ctx context.Context, page int, limit int) ([]T, *exceptions.WrappedError)
	Save(ctx context.Context, entity T) *exceptions.WrappedError
	Remove(ctx context.Context, entity T) *exceptions.WrappedError
	CorrecTimezone(entity T)
}
