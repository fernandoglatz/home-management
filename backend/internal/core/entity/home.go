package entity

type Home struct {
	Entity `bson:",inline"`
	Users  []string `json:"users,omitempty" bson:"users,omitempty"`
	Name   string   `json:"name,omitempty" bson:"name,omitempty"`
}

func (home *Home) GetCollectionName() string {
	return "homes"
}
