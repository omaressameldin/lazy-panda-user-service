package db

import (
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/badoux/checkmail"
	"github.com/golang/protobuf/ptypes"
	v1 "github.com/omaressameldin/lazy-panda-user-service/pkg/api/v1"
	"github.com/omaressameldin/lazy-panda-user-service/pkg/database"
)

func validateName(name string) error {
	if len(name) <= 3 {
		return fmt.Errorf("length should be more than 2")
	}
	return nil
}

func validateUser(email, nickname, fullname *string) []database.Validator {
	return []database.Validator{
		database.CreateValidator(
			"Email",
			func() error {
				if email != nil {
					return checkmail.ValidateFormat(*email)
				}
				return nil
			},
		),
		database.CreateValidator(
			"Nickname",
			func() error {
				if nickname != nil {
					return validateName(*nickname)
				}
				return nil
			},
		),
		database.CreateValidator(
			"Fullname",
			func() error {
				if fullname != nil {
					return validateName(*fullname)
				}
				return nil
			},
		),
	}
}

func CreateUser(connector database.Connector, key string, user *v1.User) error {
	user.CreatedAt, _ = ptypes.TimestampProto(time.Now())
	user.UpdatedAt = user.CreatedAt

	return connector.Create(
		validateUser(&user.Email, &user.Nickname, &user.Fullname),
		key,
		user,
	)
}

func ReadUser(connector database.Connector, key string) (*v1.User, error) {
	var user v1.User
	if err := connector.Read(key, &user); err != nil {
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

func UpdateUser(connector database.Connector, key string, user *v1.UserUpdate) error {
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

	return connector.Update(
		validateUser(email, nickname, fullname),
		key,
		getUpdated(user),
	)
}

func DeleteUser(connector database.Connector, key string) error {
	if err := connector.Delete(key); err != nil {
		return err
	}

	return nil
}

func ReadAllUsers(connector database.Connector) ([]*v1.User, error) {
	var users []*v1.User
	apendFn := func(i interface{}) { users = append(users, i.(*v1.User)) }
	genRefFn := func() interface{} { return &v1.User{} }

	if err := connector.ReadAll(genRefFn, apendFn); err != nil {
		return nil, err
	}

	return users, nil
}
