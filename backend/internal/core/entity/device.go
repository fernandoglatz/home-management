package entity

type Device struct {
	Entity      `bson:",inline"`
	Home        string `json:"home,omitempty" bson:"home,omitempty"`
	Name        string `json:"name,omitempty" bson:"name,omitempty"`
	Description string `json:"description,omitempty" bson:"description,omitempty"`
}

func (device Device) GetCollectionName() string {
	return "devices"
}
