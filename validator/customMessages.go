package validator

import "fmt"

func getRequiredMessage(field string) string {
	return fmt.Sprintf("%s is required!", field)
}

func getLessThanMessage(field string, lessThan string) string {
	return fmt.Sprintf("The length of %s must be less than %s!", field, lessThan)
}

func getGreaterThanMessage(field string, greaterThan string) string {
	return fmt.Sprintf("The length of %s must be greater than %s!", field, greaterThan)
}
