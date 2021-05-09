package model

type RegisterInput struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginInput struct {
	UsernameOrEmail string `json:"usernameOrEmail"`
	Password        string `json:"password"`
}

type ForgotPasswordInput struct {
	UsernameOrEmail string `json:"usernameOrEmail"`
}

type ChangePasswordInput struct {
	Token       string `json:"token"`
	NewPassword string `json:"newPassword"`
}
