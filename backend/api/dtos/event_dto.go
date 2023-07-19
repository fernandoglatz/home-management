package dtos

type EventDTO struct {
	Home         string  `json:"home,omitempty" bson:"home,omitempty" binding:"required" validate:"required"`
	Device       string  `json:"device,omitempty" bson:"device,omitempty" binding:"required" validate:"required"`
	Name         string  `json:"name,omitempty" bson:"name,omitempty" binding:"required" validate:"required"`
	Details      string  `json:"details,omitempty" bson:"details,omitempty"`
	Type         string  `json:"type,omitempty" bson:"type,omitempty" binding:"required" validate:"required"`
	TextValue    string  `json:"textValue,omitempty" bson:"textValue,omitempty"`
	NumericValue float64 `json:"numericValue,omitempty" bson:"numericValue,omitempty"`
	BooleanValue bool    `json:"booleanValue" bson:"booleanValue"`
}
