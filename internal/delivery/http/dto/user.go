package dto

type CreateUserRequest struct {
	Name  string `json:"name" binding:"required" example:"Alice"`
	Email string `json:"email" binding:"required,email" example:"alice@example.com"`
}

type UpdateUserRequest struct {
	Name  string `json:"name,omitempty" example:"Alice Updated"`
	Email string `json:"email,omitempty" binding:"omitempty,email" example:"alice.updated@example.com"`
}
