package models

import "time"

type User struct {
	Entity `bson:",inline"`
	Email  string `json:"email,omitempty" bson:"email,omitempty" binding:"required,email"`
	Name   string `json:"name,omitempty" bson:"name,omitempty" binding:"required"`
}

func (user *User) GetEntityName() string {
	return "users"
}

func (user *User) GetID() string {
	return user.ID
}

func (user *User) SetID(id string) {
	user.ID = id
}

func (user *User) GetCreatedAt() time.Time {
	return user.CreatedAt
}

func (user *User) SetCreatedAt(createdAt time.Time) {
	user.CreatedAt = createdAt
}

func (user *User) GetUpdatedAt() time.Time {
	return user.UpdatedAt
}

func (user *User) SetUpdatedAt(updatedAt time.Time) {
	user.UpdatedAt = updatedAt
}
