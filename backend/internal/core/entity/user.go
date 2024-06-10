package entity

type User struct {
	Entity `bson:",inline"`
	Email  string `json:"email,omitempty" bson:"email,omitempty"`
	Name   string `json:"name,omitempty" bson:"name,omitempty"`
}

func (user User) GetCollectionName() string {
	return "users"
}
