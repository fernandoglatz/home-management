package entity

import "time"

type Device struct {
	Entity      `bson:",inline"`
	Home        string `json:"home,omitempty" bson:"home,omitempty"`
	Name        string `json:"name,omitempty" bson:"name,omitempty"`
	Description string `json:"description,omitempty" bson:"description,omitempty"`
}

func (device *Device) GetCollectionName() string {
	return "devices"
}

func (device *Device) GetID() string {
	return device.ID
}

func (device *Device) SetID(id string) {
	device.ID = id
}

func (device *Device) GetCreatedAt() time.Time {
	return device.CreatedAt
}

func (device *Device) SetCreatedAt(createdAt time.Time) {
	device.CreatedAt = createdAt
}

func (device *Device) GetUpdatedAt() time.Time {
	return device.UpdatedAt
}

func (device *Device) SetUpdatedAt(updatedAt time.Time) {
	device.UpdatedAt = updatedAt
}
