package app

type RegisterInput struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Age       int    `json:"age"`
	IsMarried bool   `json:"is_married"`
	Password  string `json:"password"`
}
