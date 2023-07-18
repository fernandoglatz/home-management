package models

import "time"

type Entity struct {
	ID        string    `json:"id,omitempty" bson:"id,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
}

type IEntity interface {
	GetEntityName() string

	GetID() string
	SetID(string)

	GetCreatedAt() time.Time
	SetCreatedAt(time.Time)

	GetUpdatedAt() time.Time
	SetUpdatedAt(time.Time)
}
