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

type UserLoginRequest struct {
	Email string `json:"email" required:"true"`
	Password string `json:"password" required:"true"`
}