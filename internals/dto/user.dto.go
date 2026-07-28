package dto

type UserRequestType struct {
	Email    string `json:"email" required:"true"`
	Username string `json:"username" required:"true"`
	Password string `json:"password" required:"true"`
}

type UserResponseType struct {
	Email string `json:"email"`
	Username string `json:"username"`
}