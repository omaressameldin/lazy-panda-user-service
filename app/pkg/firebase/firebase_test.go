package firebase

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"
)

func getGenerateJsonErrorOutput(errs []ValidationError) error {
	j, _ := json.Marshal(errs)
	return errors.New(string(j))
}

func TestGenerateJsonError(t *testing.T) {
	errs := []ValidationError{
		{"field1", "err1"},
		{"field2", "err2"},
	}
	testTable := []struct {
		input  []ValidationError
		output error
	}{
		{[]ValidationError{}, nil},
		{errs, getGenerateJsonErrorOutput(errs)},
	}

	for _, row := range testTable {
		err := generateJsonError(row.input...)
		t.Log(err)
		if !reflect.DeepEqual(err, row.output) {
			t.Errorf("incorrect output, got: %v, want: %v.", err, row.output)
		}
	}
}

func TestCreateValidator(t *testing.T) {
	field := "test"
	err := errors.New("error message")
	fn := func() error { return err }
	validator := CreateValidator(field, fn)

	if validator.Field != field {
		t.Errorf("incorrect output, got: %v, want: %v.", validator.Field, field)
	}
	if validator.Function() != err {
		t.Errorf("incorrect output, got: %v, want: %v.", validator.Function(), err)
	}
}

func TestCombineValidationErrors(t *testing.T) {
	type validationData struct {
		f   string
		err error
	}

	validators := []Validator{}
	validationErrors := []validationData{
		{"field1", errors.New("err1")},
		{"field2", errors.New("err2")},
	}
	validationSuccesses := []validationData{
		{"field3", nil},
	}

	for _, ve := range validationErrors {
		err := ve.err
		validators = append(validators, Validator{ve.f, func() error { return err }})
	}
	for _, vs := range validationSuccesses {
		err := vs.err
		validators = append(validators, Validator{vs.f, func() error { return err }})
	}
	combined := CombineValidationErrors(validators...)

	if len(combined) != len(validationErrors) {
		t.Errorf("incorrect length, got: %d, want: %d.", len(combined), len(validationErrors))
	}
	for i, ve := range validationErrors {
		expected := ValidationError{
			Field:   ve.f,
			Message: ve.err.Error(),
		}

		if combined[i] != expected {
			t.Errorf("incorrect output, got: %s, want: %s.", combined[i], expected)
		}
	}
}
