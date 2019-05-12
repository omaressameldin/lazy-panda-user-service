package firebase

import (
	"fmt"

	"cloud.google.com/go/firestore"
	"github.com/omaressameldin/lazy-panda-user-service/pkg/database"
)

func createError(err error) error {
	return fmt.Errorf("error connecting to firebase: %v", err)
}

func addToFirebase(
	collection string,
	key string,
	validators []database.Validator,
	addFn func() error,
) error {
	errors := database.CombineValidationErrors(validators...)

	if len(errors) == 0 {
		err := addFn()
		if err != nil {
			errors = append(
				errors,
				database.ValidationError{Field: "FIREBASE", Message: err.Error()},
			)
		}
	}

	return database.GenerateJsonError(errors...)
}

func generateFirestoreUpdate(data []database.Updated) []firestore.Update {
	updated := make([]firestore.Update, 0, len(data))

	for _, item := range data {
		updated = append(updated, firestore.Update{
			Path:  item.Key,
			Value: item.Val,
		})
	}

	return updated
}
