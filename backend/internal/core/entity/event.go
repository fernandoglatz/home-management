package entity

import (
	"fernandoglatz/home-management/internal/core/entity/event"
	"time"
)

type Event struct {
	Entity  `bson:",inline"`
	Home    string     `json:"home,omitempty" bson:"home,omitempty"`
	Device  string     `json:"device,omitempty" bson:"device,omitempty"`
	Version string     `json:"version,omitempty" bson:"version,omitempty"`
	Type    event.Type `json:"type,omitempty" bson:"type,omitempty"`
	Date    time.Time  `json:"date,omitempty" bson:"date,omitempty"`
}

type IEvent interface {
	IEntity

	GetDate() time.Time
	SetDate(time.Time)
}

func (event Event) GetCollectionName() string {
	return "events"
}

func (event Event) GetDate() time.Time {
	return event.Date
}

func (event *Event) SetDate(date time.Time) {
	event.Date = date
}
