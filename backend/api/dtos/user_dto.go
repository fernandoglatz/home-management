package dtos

type UserDTO struct {
	Email string `json:"email,omitempty" bson:"email,omitempty" binding:"required,email" validate:"required"`
	Name  string `json:"name,omitempty" bson:"name,omitempty" binding:"required" validate:"required"`
}

func (userDTO *UserDTO) GetDTO() string {
	return "users"
}
