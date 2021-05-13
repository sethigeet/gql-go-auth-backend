package model

// FieldError is the structure of the error that is returned in many responses of the graphql resolvers
type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ResetPasswordResponse is the structure of the response returned by the reset password resolver
type ResetPasswordResponse struct {
	Errors     []*FieldError `json:"errors"`
	Successful bool          `json:"successful"`
}

// UserResponse is the structure of the response returned by many resolvers related to the user entity
type UserResponse struct {
	Errors []*FieldError `json:"errors"`
	User   *User         `json:"user"`
}
