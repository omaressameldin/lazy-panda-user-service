package db

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/golang/protobuf/ptypes/wrappers"
	v1 "github.com/omaressameldin/lazy-panda-user-service/pkg/api/v1"
	"github.com/stretchr/testify/mock"
	"github.com/omaressameldin/lazy-panda-user-service/pkg/database"
)

func TestValidateName(t *testing.T) {
	testString := "test"
	testTable := []struct {
		input         string
		checkOutputFn func(error)
	}{
		{
			testString,
			func(output error) {
				if output != nil {
					t.Errorf("incorrect output got: error, want: nil")
				}
			}},
		{
			testString[:MIN_NAME_LENGTH-1],
			func(output error) {
				if output == nil {
					t.Errorf("incorrect output got: nil, want: error")
				}
			}},
	}

	for _, row := range testTable {
		row.checkOutputFn(validateName(row.input))
	}
}

func TestValidateUser(t *testing.T) {
	testEmail := "test@test.test"
	testName := "test"
	emptyString := ""
	testTable := []struct {
		email          *string
		nickname       *string
		fullname       *string
		numberOfErrors int
	}{
		{
			&testEmail,
			&testName,
			&testName,
			0,
		},
		{
			&emptyString,
			&emptyString,
			&emptyString,
			3,
		},
	}

	for _, row := range testTable {
		validators := validateUser(row.email, row.nickname, row.fullname)
		c := 0
		for _, v := range validators {
			if v.Error != nil {
				c++
			}
		}
		if c != row.numberOfErrors {
			t.Errorf("incorrect output got: %v, want: %v", len(validators), row.numberOfErrors)
		}
	}
}

func TestCreateUser(t *testing.T) {
	key := "key"
	user := &v1.User{
		AuthId:   key,
		Email:    "test@test.com",
		Fullname: "test",
		Nickname: "test",
	}
	testTable := []error{
		nil,
		fmt.Errorf("test error"),
	}

	for _, expected := range testTable {
		mc := new(database.MockedConnector)
		mc.StubFnCall(
			"Create",
			validateUser(&user.Email, &user.Nickname, &user.Fullname),
			key,
			user,
		)(expected)
		err := CreateUser(mc, key, user)
		if err != expected {
			t.Errorf("incorrect output got: %v, want: %v", err, expected)
		}
	}
}

func TestReadUser(t *testing.T) {
	key := "key"
	user := &v1.User{
		AuthId:   key,
		Email:    "test@test.com",
		Fullname: "test",
		Nickname: "test",
	}
	testTable := []struct {
		key  string
		user *v1.User
		err  error
	}{
		{
			key,
			user,
			nil,
		},
		{
			key,
			nil,
			fmt.Errorf("test error"),
		},
	}

	for _, row := range testTable {
		mc := new(database.MockedConnector)
		mc.StubFnCall(
			"Read",
			key,
			&v1.User{},
		)(row.user, row.err)
		user, err := ReadUser(mc, key)
		if err != row.err {
			t.Errorf("incorrect output got: %v, want: %v", err, row.err)
		}
		if !reflect.DeepEqual(user, row.user) {
			t.Errorf("incorrect output got: %v, want: %v", user, row.user)
		}
	}
}

func TestGetUpdated(t *testing.T) {
	email := &wrappers.StringValue{Value: "test@test.test"}
	name := &wrappers.StringValue{Value: "testname"}
	pic := &wrappers.StringValue{Value: "testPic"}
	testTable := []struct {
		userUpdate *v1.UserUpdate
		expected   []database.Updated
	}{
		{
			&v1.UserUpdate{Email: email, Fullname: name, Nickname: name, Picture: pic},
			[]database.Updated{
				database.Updated{Key: "Email", Val: email.Value},
				database.Updated{Key: "Fullname", Val: name.Value},
				database.Updated{Key: "Nickname", Val: name.Value},
				database.Updated{Key: "Picture", Val: pic.Value},
			},
		},
		{
			&v1.UserUpdate{Email: email},
			[]database.Updated{database.Updated{Key: "Email", Val: email.Value}},
		},
		{
			&v1.UserUpdate{Fullname: name},
			[]database.Updated{database.Updated{Key: "Fullname", Val: name.Value}},
		},
		{
			&v1.UserUpdate{Nickname: name},
			[]database.Updated{database.Updated{Key: "Nickname", Val: name.Value}},
		},
		{
			&v1.UserUpdate{Picture: pic},
			[]database.Updated{database.Updated{Key: "Picture", Val: pic.Value}},
		},
	}

	for _, row := range testTable {
		updated := getUpdated(row.userUpdate)
		if !reflect.DeepEqual(updated, row.expected) {
			t.Errorf("incorrect output got: %v, want: %v", updated, row.expected)
		}
	}
}

func TestUpdateUser(t *testing.T) {
	key := "key"
	user := &v1.UserUpdate{
		Email:    &wrappers.StringValue{Value: "test@test.test"},
		Fullname: &wrappers.StringValue{Value: "test"},
		Nickname: &wrappers.StringValue{Value: "test"},
	}
	testTable := []error{
		nil,
		fmt.Errorf("test error"),
	}

	for _, expected := range testTable {
		mc := new(database.MockedConnector)
		mc.StubFnCall(
			"Update",
			validateUser(&user.Email.Value, &user.Nickname.Value, &user.Fullname.Value),
			key,
			getUpdated(user),
		)(expected)
		err := UpdateUser(mc, key, user)
		if err != expected {
			t.Errorf("incorrect output got: %v, want: %v", err, expected)
		}
	}
}

func TestDeleteUser(t *testing.T) {
	key := "key"
	testTable := []error{
		nil,
		fmt.Errorf("test error"),
	}

	for _, expected := range testTable {
		mc := new(database.MockedConnector)
		mc.StubFnCall(
			"Delete",
			key,
		)(expected)
		err := DeleteUser(mc, key)
		if err != expected {
			t.Errorf("incorrect output got: %v, want: %v", err, expected)
		}
	}
}

func TestReadAllUsers(t *testing.T) {
	key := "key"
	key2 := "key2"
	user := &v1.User{
		AuthId:   key,
		Email:    "test@test.com",
		Fullname: "test",
		Nickname: "test",
	}
	user2 := &v1.User{
		AuthId:   key2,
		Email:    "tes2t@test.com",
		Fullname: "test2",
		Nickname: "test2",
	}
	testTable := []struct {
		key  string
		users []*v1.User
		err  error
	}{
		{
			key,
			[]*v1.User{user, user2},
			nil,
		},
		{
			key,
			nil,
			fmt.Errorf("test error"),
		},
	}

	for _, row := range testTable {
		mc := new(database.MockedConnector)
		mc.StubFnCall(
			"ReadAll",
			// AnythingOfType does not work with interfaces
			// source: https://github.com/stretchr/testify/issues/68
			mock.Anything,
			mock.Anything,
		)(row.users, row.err)
		users, err := ReadAllUsers(mc)
		if err != row.err {
			t.Errorf("incorrect output got: %v, want: %v", err, row.err)
		}
		if !reflect.DeepEqual(users, row.users) {
			t.Errorf("incorrect output got: %v, want: %v", users, row.users)
		}
	}
}
