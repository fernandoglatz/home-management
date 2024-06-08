package entity

type RFEvent struct {
	Event            `bson:",inline"`
	Frequency        int `json:"frequency" bson:"frequency"`
	Code             int `json:"code" bson:"code"`
	Bits             int `json:"bits" bson:"bits"`
	Protocol         int `json:"protocol" bson:"protocol"`
	ReceiveTolerance int `json:"receiveTolerance" bson:"receiveTolerance"`
}
