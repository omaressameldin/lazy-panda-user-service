package firebase

import (
	"encoding/json"
	"fmt"

	"cloud.google.com/go/firestore"
	"golang.org/x/net/context"
	"google.golang.org/api/iterator"
)

type ValidationError struct {
	Field   string
	Message string
}

type Validator struct {
	Field    string
	Function func() error
}

func generateJsonError(errors ...ValidationError) error {
	if len(errors) > 0 {
		jsonError, err := json.Marshal(errors)
		if err != nil {
			return err
		}
		return fmt.Errorf(string(jsonError))
	}
	return nil
}

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

func Create(collection string, key string, data interface{}, validationFn func() []ValidationError) error {
	errors := validationFn()
	if len(errors) == 0 {
		_, err := client.Collection(collection).Doc(key).Set(context.Background(), data)
		if err != nil {
			errors = append(errors, ValidationError{Field: "FIREBASE", Message: err.Error()})
		}
	}
	return generateJsonError(errors...)
}

func Update(collection string, key string, data []firestore.Update, validationFn func() []ValidationError) error {
	errors := validationFn()
	if len(errors) == 0 {
		_, err := client.Collection(collection).Doc(key).Update(context.Background(), data)
		if err != nil {
			errors = append(errors, ValidationError{Field: "FIREBASE", Message: err.Error()})
		}
	}
	return generateJsonError(errors...)
}

func Read(collection string, key string, model interface{}) error {
	var err error
	docSnap, err := client.Collection(collection).Doc(key).Get(context.Background())
	if err == nil {
		err = docSnap.DataTo(model)
	}
	if err != nil {
		return generateJsonError(ValidationError{Field: "FIREBASE", Message: err.Error()})
	}

	return nil
}

func Delete(collection string, key string) error {
	if _, err := client.Collection(collection).Doc(key).Delete(context.Background()); err != nil {
		return generateJsonError(ValidationError{Field: "FIREBASE", Message: err.Error()})
	}

	return nil
}

func ReadAll(collection string, genRefFn func() interface{}, appendFn func(interface{})) error {
	docs := client.Collection(collection).Documents(context.Background())

	for {
		docSnap, err := docs.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return generateJsonError(ValidationError{Field: "FIREBASE", Message: err.Error()})
		}

		recordRef := genRefFn()
		err = docSnap.DataTo(recordRef)
		if err != nil {
			return generateJsonError(ValidationError{Field: "FIREBASE", Message: err.Error()})
		}
		appendFn(recordRef)
	}
	return nil
}
