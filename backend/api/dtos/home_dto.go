package dtos

type HomeDTO struct {
	Name  string   `json:"name,omitempty" bson:"name,omitempty" binding:"required" validate:"required"`
	Users []string `json:"users,omitempty" bson:"users,omitempty"`
}
