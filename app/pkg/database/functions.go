package database

import (
	"encoding/json"
	"fmt"
)

func CreateValidator(field string, function func() error) Validator {
	return Validator{
		Field:    field,
		Function: function,
	}
}

func CombineValidationErrors(validators ...Validator) []ValidationError {
	combinedErrors := make([]ValidationError, 0, len(validators))
	for _, v := range validators {
		if err := v.Function(); err != nil {
			combinedErrors = append(
				combinedErrors,
				ValidationError{
					Field:   v.Field,
					Message: err.Error(),
				},
			)
		}
	}
	return combinedErrors
}

func GenerateJsonError(errors ...ValidationError) error {
	if len(errors) > 0 {
		jsonError, _ := json.Marshal(errors)

		return fmt.Errorf(string(jsonError))
	}
	return nil
}
