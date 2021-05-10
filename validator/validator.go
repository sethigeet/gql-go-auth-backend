// Package validator provides a function to validate structs according to the
// validations provided in the struct tags
package validator

import (
	"strings"

	"github.com/go-playground/validator"
	"github.com/sethigeet/gql-go-auth-backend/graph/model"
)

func Validate(toValidate interface{}) []*model.FieldError {
	validate := validator.New()
	err := validate.Struct(toValidate)

	if err != nil {
		errs := []*model.FieldError{}
		for _, err := range err.(validator.ValidationErrors) {
			errs = append(errs, translate(err))
		}
		return errs
	}

	return nil
}

func translate(err validator.FieldError) *model.FieldError {
	field := err.Field()
	errTag := err.Tag()

	var message string
	switch errTag {
	case "required":
		message = getRequiredMessage(field)
	case "email":
		message = getEmailMessage(field)
	case "excludesrune":
		message = getExcludesMessage(field, err.Param())
	case "lt":
		message = getLessThanMessage(field, err.Param())
	case "gt":
		message = getGreaterThanMessage(field, err.Param())
	default:
		panic("This type of validation error is not implemented yet!")
	}

	return &model.FieldError{
		Field:   strings.ToLower(string(field[0])) + string(field[1:]),
		Message: message,
	}
}
