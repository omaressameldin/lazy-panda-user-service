package db

import (
	"fmt"
	"time"

	"github.com/badoux/checkmail"
	"github.com/golang/protobuf/ptypes"
	v1 "github.com/omaressameldin/lazy-panda-user-service/app/pkg/api/v1"
	"github.com/omaressameldin/lazy-panda-utils/app/pkg/database"
)

const MIN_NAME_LENGTH int = 3

func validateName(name string) error {
	if len(name) <= MIN_NAME_LENGTH {
		return fmt.Errorf("length should be at least %d", MIN_NAME_LENGTH)
	}
	return nil
}

func validateUser(email, nickname, fullname *string) []database.Validator {
	var emailError error
	if email != nil {
		emailError = checkmail.ValidateFormat(*email)
	}
	var nicknameError error
	if nickname != nil {
		nicknameError = validateName(*nickname)
	}
	var fullnameError error
	if nickname != nil {
		fullnameError = validateName(*fullname)
	}

	return []database.Validator{
		database.CreateValidator(
			"Email",
			emailError,
		),
		database.CreateValidator(
			"Nickname",
			nicknameError,
		),
		database.CreateValidator(
			"Fullname",
			fullnameError,
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

func getUpdated(user *v1.UserUpdate) []database.Updated {
	updated := []database.Updated{}
	if user.Email != nil {
		updated = append(updated, database.Updated{Key: "Email", Val: user.Email.Value})
	}
	if user.Fullname != nil {
		updated = append(updated, database.Updated{Key: "Fullname", Val: user.Fullname.Value})
	}
	if user.Nickname != nil {
		updated = append(updated, database.Updated{Key: "Nickname", Val: user.Nickname.Value})
	}
	if user.Picture != nil {
		updated = append(updated, database.Updated{Key: "Picture", Val: user.Picture.Value})
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
