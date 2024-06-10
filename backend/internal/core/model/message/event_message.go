package message

import (
	"fernandoglatz/home-management/internal/core/entity/event"
	"time"
)

type EventMessage struct {
	Type    event.Type `json:"type,omitempty"`
	Device  string     `json:"device,omitempty"`
	Version string     `json:"version,omitempty"`
	Date    time.Time  `json:"date,omitempty"`
}
