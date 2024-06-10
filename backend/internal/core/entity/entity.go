package entity

import "time"

type Entity struct {
	ID        string    `json:"id,omitempty" bson:"id,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
}

func (entity Entity) GetID() string {
	return entity.ID
}

func (entity *Entity) SetID(id string) {
	entity.ID = id
}

func (entity Entity) GetCreatedAt() time.Time {
	return entity.CreatedAt
}

func (entity *Entity) SetCreatedAt(createdAt time.Time) {
	entity.CreatedAt = createdAt
}

func (entity Entity) GetUpdatedAt() time.Time {
	return entity.UpdatedAt
}

func (entity *Entity) SetUpdatedAt(updatedAt time.Time) {
	entity.UpdatedAt = updatedAt
}

type IEntity interface {
	GetCollectionName() string

	GetID() string
	SetID(string)

	GetCreatedAt() time.Time
	SetCreatedAt(time.Time)

	GetUpdatedAt() time.Time
	SetUpdatedAt(time.Time)
}
