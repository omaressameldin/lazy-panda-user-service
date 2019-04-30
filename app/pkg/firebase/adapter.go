package firebase

import (
	"fmt"

	"github.com/omaressameldin/lazy-panda-user-service/pkg/database"
	"golang.org/x/net/context"
	"google.golang.org/api/iterator"
)

func (f *Firebase) Create(
	validators []database.Validator,
	key string,
	data interface{},
) error {
	return addToFirebase(
		f.Collection,
		key,
		validators,
		func() error {
			_, err := f.client.Collection(f.Collection).Doc(key).Set(
				context.Background(),
				data,
			)
			return err
		},
	)
}

func (f *Firebase) Update(
	validators []database.Validator,
	key string,
	data []database.Updated,
) error {
	fmt.Println(generateFirestoreUpdate(data))
	return addToFirebase(
		f.Collection,
		key,
		validators,
		func() error {
			_, err := f.client.Collection(f.Collection).Doc(key).Update(
				context.Background(),
				generateFirestoreUpdate(data),
			)
			return err
		},
	)
}

func (f *Firebase) Read(key string, model interface{}) error {
	var err error
	docSnap, err := f.client.Collection(f.Collection).Doc(key).Get(context.Background())
	if err == nil {
		err = docSnap.DataTo(model)
	}
	if err != nil {
		return database.GenerateJsonError(database.ValidationError{
			Field:   "FIREBASE",
			Message: err.Error(),
		})
	}

	return nil
}

func (f *Firebase) Delete(key string) error {
	_, err := f.client.Collection(f.Collection).Doc(key).Delete(context.Background())
	if err != nil {
		return database.GenerateJsonError(database.ValidationError{
			Field:   "FIREBASE",
			Message: err.Error(),
		})
	}

	return nil
}

func (f *Firebase) ReadAll(
	genRefFn func() interface{},
	appendFn func(interface{}),
) error {
	docs := f.client.Collection(f.Collection).Documents(context.Background())

	for {
		docSnap, err := docs.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return database.GenerateJsonError(database.ValidationError{
				Field:   "FIREBASE",
				Message: err.Error(),
			})
		}

		recordRef := genRefFn()
		err = docSnap.DataTo(recordRef)
		if err != nil {
			return database.GenerateJsonError(database.ValidationError{
				Field:   "FIREBASE",
				Message: err.Error(),
			})
		}
		appendFn(recordRef)
	}
	return nil
}
