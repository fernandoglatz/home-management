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

func GetType(typeStr string) Type {
	types := make(map[string]Type)
	types[string(RECEIVED_RF)] = RECEIVED_RF
	types[string(SEND_RF)] = SEND_RF
	types[string(MQTT_CONNECTED)] = MQTT_CONNECTED
	types[string(RESTART)] = RESTART
	types[string(RESET)] = RESET
	types[string(SET_CONFIG)] = SET_CONFIG
	types[string(GET_CONFIG)] = GET_CONFIG
	types[string(GET_INFO)] = GET_INFO
	types[string(ACTION_UNRECOGNIZED)] = ACTION_UNRECOGNIZED

	return types[typeStr]
}
