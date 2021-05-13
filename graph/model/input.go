package model

// RegisterInput is the structure input that is received by the register resolver
type RegisterInput struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required,excludesrune=@"`
	Password string `json:"password" validate:"required,gt=4"`
}

// LoginInput is the structure input that is received by the login resolver
type LoginInput struct {
	UsernameOrEmail string `json:"usernameOrEmail" validate:"required"`
	Password        string `json:"password" validate:"required"`
}

// ForgotPasswordInput is the structure input that is received by the forgot password resolver
type ForgotPasswordInput struct {
	UsernameOrEmail string `json:"usernameOrEmail" validate:"required"`
}

// ChangePasswordInput is the structure input that is received by the change password resolver
type ChangePasswordInput struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"newPassword" validate:"required,gt=4"`
}
