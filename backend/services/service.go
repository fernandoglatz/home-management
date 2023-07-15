package services

import (
	"github.com/fernandoglatz/home-management/models"
	"github.com/fernandoglatz/home-management/repositories"
)

type service[T models.IEntity] struct {
	baseEntity T
	repository *repositories.Repository[T]
}

func NewService[T models.IEntity]() *service[T] {
	service := &service[T]{}
	service.repository = repositories.NewRepository[T]()
	return service
}
