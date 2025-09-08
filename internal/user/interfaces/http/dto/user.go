package dto

type RegisterRequest struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Age       int    `json:"age" validate:"required,min=18"`
	Email     string `json:"email" validate:"required,email"`
	IsMarried bool   `json:"is_married"`
	Password  string `json:"password" validate:"required,min=8"`
}

type UserResponse struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	FullName  string `json:"full_name"`
	Age       int    `json:"age"`
	IsMarried bool   `json:"is_married"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
