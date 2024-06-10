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

func (event *Event) GetCollectionName() string {
	return "events"
}
