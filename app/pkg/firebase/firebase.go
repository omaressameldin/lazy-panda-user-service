package firebase

import (
	"fmt"

	"cloud.google.com/go/firestore"
	"golang.org/x/net/context"

	firebase "firebase.google.com/go"

	"google.golang.org/api/option"
)

var client *firestore.Client

func createError(err error) error {
	return fmt.Errorf("error initializing app: %v", err)
}

func StartConnection(jsonConfig string) error {
	opt := option.WithCredentialsFile(jsonConfig)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return createError(err)
	}

	if client, err = app.Firestore(context.Background()); err != nil {
		return createError(err)
	}

	return nil
}

func CloseConnection() error {
	return client.Close()
}
