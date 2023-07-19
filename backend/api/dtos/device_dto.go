package dtos

type DeviceDTO struct {
	Home        string `json:"home,omitempty" bson:"home,omitempty" binding:"required" validate:"required"`
	Name        string `json:"name,omitempty" bson:"name,omitempty" binding:"required" validate:"required"`
	Description string `json:"description,omitempty" bson:"description,omitempty"`
}
