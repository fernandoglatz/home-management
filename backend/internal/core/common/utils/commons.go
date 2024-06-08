package utils

import (
	"fernandoglatz/home-management/internal/core/common/utils/constants"
	"os"
	"strings"
)

func IsEmptyStr(value string) bool {
	return len(value) == constants.ZERO
}

func IsNotEmptyStr(value string) bool {
	return len(value) > constants.ZERO
}

func IsBlankStr(value string) bool {
	return IsEmptyStr(strings.TrimSpace(value))
}

func IsNotBlankStr(value string) bool {
	return IsNotEmptyStr(strings.TrimSpace(value))
}

func GetTimezone() string {
	return os.Getenv("TZ")
}
