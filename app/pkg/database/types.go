package database

type ValidationError struct {
	Field   string
	Message string
}

type Validator struct {
	Field string
	Error error
}

type Updated struct {
	Key string
	Val interface{}
}

type Connector interface {
	CloseConnection() error
	Create(validators []Validator, key string, data interface{}) error
	Update(validators []Validator, key string, data []Updated) error
	Read(key string, model interface{}) error
	ReadAll(genRefFn func() interface{}, appendFn func(interface{})) error
	Delete(key string) error
}
