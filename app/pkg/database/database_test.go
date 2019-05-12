package database

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
		err := GenerateJsonError(row.input...)
		if !reflect.DeepEqual(err, row.output) {
			t.Errorf("incorrect output, got: %v, want: %v.", err, row.output)
		}
	}
}

func TestCreateValidator(t *testing.T) {
	field := "test"
	err := errors.New("error message")
	validator := CreateValidator(field, err)

	if validator.Field != field {
		t.Errorf("incorrect output, got: %v, want: %v.", validator.Field, field)
	}
	if validator.Error != err {
		t.Errorf("incorrect output, got: %v, want: %v.", validator.Error, err)
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
		validators = append(validators, Validator{ve.f, err})
	}
	for _, vs := range validationSuccesses {
		err := vs.err
		validators = append(validators, Validator{vs.f, err})
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
