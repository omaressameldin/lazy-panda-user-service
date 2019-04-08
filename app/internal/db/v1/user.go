package db

import (
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/badoux/checkmail"
	"github.com/golang/protobuf/ptypes"
	v1 "github.com/omaressameldin/lazy-panda-user-service/pkg/api/v1"
	"github.com/omaressameldin/lazy-panda-user-service/pkg/firebase"
)

func validateName(name string) error {
	if len(name) <= 3 {
		return fmt.Errorf("length should be more than 2")
	}
	return nil
}

func validateUser(email, nickname, fullname *string) func() []firebase.ValidationError {
	return func() []firebase.ValidationError {
		return firebase.CombineValidationErrors(
			firebase.CreateValidator(
				"Email",
				func() error {
					if email != nil {
						return checkmail.ValidateFormat(*email)
					}
					return nil
				},
			),
			firebase.CreateValidator(
				"Nickname",
				func() error {
					if nickname != nil {
						return validateName(*nickname)
					}
					return nil
				},
			),
			firebase.CreateValidator(
				"Fullname",
				func() error {
					if fullname != nil {
						return validateName(*fullname)
					}
					return nil
				},
			),
		)
	}
}

func CreateUser(collection string, key string, user *v1.User) error {
	user.CreatedAt, _ = ptypes.TimestampProto(time.Now())
	user.UpdatedAt = user.CreatedAt

	return firebase.Create(
		collection,
		key,
		user,
		validateUser(&user.Email, &user.Nickname, &user.Fullname))
}

func ReadUser(collection string, key string) (*v1.User, error) {
	var user v1.User
	if err := firebase.Read(collection, key, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func getUpdated(user *v1.UserUpdate) []firestore.Update {
	updated := []firestore.Update{}
	if user.Email != nil {
		updated = append(updated, firestore.Update{Path: "Email", Value: user.Email.Value})
	}
	if user.Fullname != nil {
		updated = append(updated, firestore.Update{Path: "Fullname", Value: user.Fullname.Value})
	}
	if user.Nickname != nil {
		updated = append(updated, firestore.Update{Path: "Nickname", Value: user.Nickname.Value})
	}
	if user.Picture != nil {
		updated = append(updated, firestore.Update{Path: "Picture", Value: user.Picture.Value})
	}
	user.UpdatedAt, _ = ptypes.TimestampProto(time.Now())
	return updated
}

func UpdateUser(collection string, key string, user *v1.UserUpdate) error {
	var email *string
	var nickname *string
	var fullname *string
	if user.Email != nil {
		email = &user.Email.Value
	}
	if user.Nickname != nil {
		nickname = &user.Email.Value
	}
	if user.Fullname != nil {
		fullname = &user.Fullname.Value
	}

	return firebase.Update(
		collection,
		key,
		getUpdated(user),
		validateUser(email, nickname, fullname),
	)
}

func DeleteUser(collection string, key string) error {
	if err := firebase.Delete(collection, key); err != nil {
		return err
	}

	return nil
}

func ReadAllUsers(collection string) ([]*v1.User, error) {
	var users []*v1.User
	apendFn := func(i interface{}) { users = append(users, i.(*v1.User)) }
	genRefFn := func() interface{} { return &v1.User{} }

	if err := firebase.ReadAll(collection, genRefFn, apendFn); err != nil {
		return nil, err
	}

	return users, nil
}
