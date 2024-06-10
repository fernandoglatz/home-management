package message

type RfEventMessage struct {
	EventMessage
	Code             int `json:"code"`
	Bits             int `json:"bits"`
	Protocol         int `json:"protocol"`
	Frequency        int `json:"frequency"`
	ReceiveTolerance int `json:"receiveTolerance"`
}
