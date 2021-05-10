package validator

import "fmt"

func getRequiredMessage(field string) string {
	return fmt.Sprintf("%s is required!", field)
}

func getEmailMessage(field string) string {
	return fmt.Sprintf("%s must be a valid email address!", field)
}

func getExcludesMessage(field string, excludes string) string {
	return fmt.Sprintf("%s must not contain '%s'", field, excludes)
}

func getLessThanMessage(field string, lessThan string) string {
	return fmt.Sprintf("The length of %s must be less than %s!", field, lessThan)
}

func getGreaterThanMessage(field string, greaterThan string) string {
	return fmt.Sprintf("The length of %s must be greater than %s!", field, greaterThan)
}
