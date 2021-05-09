package model

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ResetPasswordResponse struct {
	Errors     []*FieldError `json:"errors"`
	Successful *bool         `json:"successful"`
}

type UserResponse struct {
	Errors []*FieldError `json:"errors"`
	User   *User         `json:"user"`
}
