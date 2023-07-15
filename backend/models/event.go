package models

type Event struct {
	Entity      `bson:",inline"`
	Home        string  `json:"home,omitempty" bson:"home,omitempty" binding:"required"`
	Name        string  `json:"name,omitempty" bson:"name,omitempty" binding:"required"`
	Description string  `json:"description,omitempty" bson:"description,omitempty"`
	Type        string  `json:"type,omitempty" bson:"type,omitempty" binding:"required"`
	Value       float64 `json:"value,omitempty" bson:"value,omitempty"`
	State       bool
}

func (event Event) GetCollectionName() string {
	return "events"
}
