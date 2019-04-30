package firebase

import (
	"cloud.google.com/go/firestore"
	"golang.org/x/net/context"

	firebase "firebase.google.com/go"

	"google.golang.org/api/option"
)

type Firebase struct {
	Collection string
	client     *firestore.Client
}

func StartConnection(jsonConfig string, collection string) (*Firebase, error) {
	opt := option.WithCredentialsFile(jsonConfig)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, createError(err)
	}

	client, err := app.Firestore(context.Background())
	if err != nil {
		return nil, createError(err)
	}

	return &Firebase{
		collection,
		client,
	}, nil
}

func (f *Firebase) CloseConnection() error {
	return f.client.Close()
}
