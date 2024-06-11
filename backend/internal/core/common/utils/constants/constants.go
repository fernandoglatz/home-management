package constants

type ContextKey string

const (
	TRACE_MAP ContextKey = "TRACE-MAP"

	LOGGING_LEVEL = "LOGGING_LEVEL"
	PROFILE       = "PROFILE"
	DEV_PROFILE   = "dev"

	ID         string = "id"
	REQUEST_ID string = "REQUEST-ID"
	MESSAGE_ID string = "MESSAGE-ID"

	MINUS_ONE = -1
	ZERO      = 0
	ONE       = 1
	TEN       = 10

	EMPTY    = ""
	SLASH    = "/"
	DOT      = "."
	PLUS     = "+"
	ASTERISK = "*"
	HASH     = "#"
	COLON    = ":"
	HYPHEN   = "-"
)
