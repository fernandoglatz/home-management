package models

import "time"

type Home struct {
	Entity `bson:",inline"`
	Users  []string `json:"users,omitempty" bson:"users,omitempty"`
	Name   string   `json:"name,omitempty" bson:"name,omitempty"`
}

func (home *Home) GetCollectionName() string {
	return "homes"
}

func (home *Home) GetID() string {
	return home.ID
}

func (home *Home) SetID(id string) {
	home.ID = id
}

func (home *Home) GetCreatedAt() time.Time {
	return home.CreatedAt
}

func (home *Home) SetCreatedAt(createdAt time.Time) {
	home.CreatedAt = createdAt
}

func (home *Home) GetUpdatedAt() time.Time {
	return home.UpdatedAt
}

func (home *Home) SetUpdatedAt(updatedAt time.Time) {
	home.UpdatedAt = updatedAt
}
