package entity

import "time"

type Event struct {
	Entity  `bson:",inline"`
	Home    string `json:"home,omitempty" bson:"home,omitempty"`
	Device  string `json:"device,omitempty" bson:"device,omitempty"`
	Type    string `json:"type,omitempty" bson:"type,omitempty"`
	Version string `json:"version,omitempty" bson:"version,omitempty"`
}

func (event *Event) GetCollectionName() string {
	return "events"
}

func (event *Event) GetID() string {
	return event.ID
}

func (event *Event) SetID(id string) {
	event.ID = id
}

func (event *Event) GetCreatedAt() time.Time {
	return event.CreatedAt
}

func (event *Event) SetCreatedAt(createdAt time.Time) {
	event.CreatedAt = createdAt
}

func (event *Event) GetUpdatedAt() time.Time {
	return event.UpdatedAt
}

func (event *Event) SetUpdatedAt(updatedAt time.Time) {
	event.UpdatedAt = updatedAt
}
