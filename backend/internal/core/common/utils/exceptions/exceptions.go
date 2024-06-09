package exceptions

import "fernandoglatz/home-management/internal/core/common/utils/constants"

var (
	GenericError = BaseError{
		Code:    "GENERIC_ERROR",
		Message: "Try again later.",
	}
	RecordNotFound = BaseError{
		Code:    "RECORD_NOT_FOUND",
		Message: "Record not found.",
	}
	HeaderNotFound = BaseError{
		Code:    "HEADER_NOT_FOUND",
		Message: "Header not found.",
	}
	InvalidJSON = BaseError{
		Code:    "INVALID_JSON",
		Message: "Invalid JSON.",
	}
	DuplicatedRecord = BaseError{
		Code:    "DUPLICATED_RECORD",
		Message: "Record duplicated.",
	}
	WakeOnLanError = BaseError{
		Code:    "WAKE_ON_LAN_ERROR",
		Message: "Erro on waking up device",
	}
)

type WrappedError struct {
	Error     error
	Message   string
	Code      string
	BaseError BaseError
}

type BaseError struct {
	Code    string
	Message string
}

type ApiError struct {
	Message string
	Status  int
}

func (wrappedError WrappedError) GetMessage() string {
	if wrappedError.Error != nil {
		return wrappedError.Error.Error()
	}

	if wrappedError.Message != constants.EMPTY {
		return wrappedError.Message
	}

	if wrappedError.BaseError.Message != constants.EMPTY {
		return wrappedError.BaseError.Message
	}

	return constants.EMPTY
}

func (wrappedError WrappedError) GetCode() string {
	if wrappedError.Code != constants.EMPTY {
		return wrappedError.Code
	}

	if wrappedError.BaseError.Code != constants.EMPTY {
		return wrappedError.BaseError.Code
	}

	return constants.EMPTY
}
