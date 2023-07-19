package dtos

type ResponseDTO struct {
	Message string `json:"message,omitempty"`
	Details string `json:"details,omitempty"`
}
