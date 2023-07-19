package services

import (
	"time"

	"github.com/fernandoglatz/home-management/backend/models"
	"github.com/fernandoglatz/home-management/backend/repositories"
)

type Service[T models.IEntity] struct {
	BaseEntity T
	repository repositories.Repository[T]
}

func NewService[T models.IEntity]() Service[T] {
	return Service[T]{
		repository: repositories.NewRepository[T](),
	}
}

func (service *Service[T]) Save(entity T) error {
	var err error

	repository := service.repository
	id := entity.GetID()

	now := time.Now()
	entity.SetUpdatedAt(now)

	if entity.GetCreatedAt().IsZero() {
		entity.SetCreatedAt(now)
	}

	if id == "" {
		err = repository.Insert(entity)
	} else {
		err = repository.Update(entity)
	}

	return err
}

func (service *Service[T]) Delete(entity T) error {
	repository := service.repository
	return repository.Delete(entity)
}

func (service *Service[T]) FindByID(id string) (T, error) {
	repository := service.repository
	return repository.FindByID(id)
}

func (service *Service[T]) FindAll() ([]T, error) {
	repository := service.repository
	return repository.FindAll()
}
