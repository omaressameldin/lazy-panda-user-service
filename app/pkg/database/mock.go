package database

import (
	"reflect"
	"github.com/stretchr/testify/mock"
)

type MockedConnector struct {
	mock.Mock
}

func (mc *MockedConnector) CloseConnection() error {
	args := mc.Called()

	return args.Error(0)
}

func (mc *MockedConnector) Create(
	validators []Validator,
	key string,
	data interface{},
) error {
	args := mc.Called(validators, key, data)
	return args.Error(0)
}

func (mc *MockedConnector) Update(
	validators []Validator,
	key string,
	data []Updated,
) error {
	args := mc.Called(validators, key, data)

	return args.Error(0)
}

func (mc *MockedConnector) Read(key string, model interface{}) error {
	args := mc.Called(key, model)
	v := reflect.ValueOf(model)
	modelRef := v.Elem()
	value := reflect.ValueOf(args.Get(0))
	if !value.IsNil() {
		modelRef.Set(value.Elem())
	}
	return args.Error(1)
}

func (mc *MockedConnector) ReadAll(
	genRefFn func() interface{},
	appendFn func(interface{}),
) error {
	args := mc.Called(genRefFn, appendFn)
	slice := reflect.ValueOf(args.Get(0))
	if args.Error(1) == nil {
		for i := 0 ; i < slice.Len(); i++ {
			element := slice.Index(i).Elem()
			ref := genRefFn()
			refValue := reflect.ValueOf(ref).Elem()
			refValue.Set(element)
			appendFn(ref)
		}
	}
	return args.Error(1)
}

func (mc *MockedConnector) Delete(key string) error {
	args := mc.Called(key)

	return args.Error(0)
}

func NewMockedConnector() *MockedConnector {
	return new(MockedConnector)
}

func (mc *MockedConnector) StubFnCall(
	fn string,
	args ...interface{},
) func(...interface{}) {
	return func(toBeReturned ...interface{}) { mc.On(fn, args...).Return(toBeReturned...) }
}
