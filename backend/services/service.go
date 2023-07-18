package services

import (
	"context"
	"time"

	"github.com/fernandoglatz/home-management/models"
	"github.com/fernandoglatz/home-management/repositories"
)

type Service[T models.IEntity] struct {
	baseEntity T
	repository *repositories.Repository[T]
}

func NewService[T models.IEntity]() *Service[T] {
	service := &Service[T]{}
	service.repository = repositories.NewRepository[T]()
	return service
}

func (service *Service[T]) Save(ctx context.Context, entity T) error {
	var err error

	repository := service.repository
	id := entity.GetID()

	now := time.Now()
	entity.SetUpdatedAt(now)

	if entity.GetCreatedAt().IsZero() {
		entity.SetCreatedAt(now)
	}

	if id == "" {
		err = repository.Insert(ctx, entity)
	} else {
		err = repository.Update(ctx, entity)
	}

	return err
}

func (service *Service[T]) Delete(ctx context.Context, entity T) error {
	repository := service.repository
	return repository.Delete(ctx, entity)
}

func (service *Service[T]) FindByID(ctx context.Context, id string) (T, error) {
	repository := service.repository
	return repository.FindByID(ctx, id)
}

func (service *Service[T]) FindAll(ctx context.Context, id string) ([]T, error) {
	repository := service.repository
	return repository.FindAll(ctx)
}
