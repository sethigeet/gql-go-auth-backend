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

// GetAlreadyExistsMessage returns a string that can be returned back to the user
// as a feedback message to tell them that an entry with that value already exists
func GetAlreadyExistsMessage(field string) string {
	return fmt.Sprintf("An entry with that %s already exists!", field)
}
