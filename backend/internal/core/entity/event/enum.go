package event

type Type string

const (
	RECEIVED_RF         Type = "RECEIVED_RF"
	SEND_RF             Type = "SEND_RF"
	MQTT_CONNECTED      Type = "MQTT_CONNECTED"
	RESTART             Type = "RESTART"
	RESET               Type = "RESET"
	SET_CONFIG          Type = "SET_CONFIG"
	GET_CONFIG          Type = "GET_CONFIG"
	GET_INFO            Type = "GET_INFO"
	ACTION_UNRECOGNIZED Type = "ACTION_UNRECOGNIZED"
)
