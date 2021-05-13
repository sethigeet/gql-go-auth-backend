package validator

import "fmt"

func getRequiredMessage(field string) string {
	return fmt.Sprintf("%s is required!", field)
}

func getInvalidEmailMessage(field string) string {
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

// GetDoesNotExistMessage returns a string that can be returned back to the user
// as a feedback message to tell them that an entry with that value does not exist
func GetDoesNotExistMessage(field string) string {
	return fmt.Sprintf("An entry with that %s does not exist!", field)
}

// GetIncorrectMessage returns a string that can be returned back to the user
// as a feedback message to tell them that the value they entered is incorrect
func GetIncorrectMessage(field string) string {
	return fmt.Sprintf("The %s you entered is incorrect!", field)
}

func getInvalidUUIDMessage(field string) string {
	return fmt.Sprintf("The %s you entered is not a valid uuid!", field)
}

// GetInvalidTokenMessage returns a string that can be returned back to the user
// as a feedback message to tell them that the value they entered is not a valid token
func GetInvalidTokenMessage(field string) string {
	return fmt.Sprintf("The %s is either invalid or has expired!", field)
}
