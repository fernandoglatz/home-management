package exceptions

type BusinessException struct {
	Message string
	Details string
}

func (e *BusinessException) Error() string {
	return e.Message
}
