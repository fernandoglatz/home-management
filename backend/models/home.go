package models

type Home struct {
	Entity `bson:",inline"`
	User   string `json:"user,omitempty" bson:"user,omitempty" binding:"required"`
	Name   string `json:"name,omitempty" bson:"name,omitempty" binding:"required"`
}

func (home Home) GetCollectionName() string {
	return "homes"
}
