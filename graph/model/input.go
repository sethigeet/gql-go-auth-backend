package model

type RegisterInput struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required,excludesrune=@"`
	Password string `json:"password" validate:"required,gt=4"`
}

type LoginInput struct {
	UsernameOrEmail string `json:"usernameOrEmail" validate:"required"`
	Password        string `json:"password" validate:"required"`
}

type ForgotPasswordInput struct {
	UsernameOrEmail string `json:"usernameOrEmail" validate:"required"`
}

type ChangePasswordInput struct {
	Token       string `json:"token" validate:"required,uuid4"`
	NewPassword string `json:"newPassword" validate:"required,gt=4"`
}
