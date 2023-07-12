package models

type User struct {
	ID    string `json:"id,omitempty" bson:"id,omitempty"`
	Email string `json:"email,omitempty" bson:"email,omitempty" binding:"required"`
	Name  string `json:"name,omitempty" bson:"name,omitempty" binding:"required"`
}
