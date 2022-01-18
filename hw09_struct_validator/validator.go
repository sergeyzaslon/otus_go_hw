package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var (
	ErrStrLen       = errors.New("string length does not match expected")
	ErrNumRange     = errors.New("number is out of range")
	ErrInvalidEmail = errors.New("invalid email")
	ErrStrEnum      = errors.New("the value is not in allowed enum")
	ErrRegexp       = errors.New("the value does not match regexp pattern")
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	messages := make([]string, 0, len(v))
	for _, err := range v {
		messages = append(messages, err.Err.Error())
	}
	return strings.Join(messages, "\n")
}

func Validate(v interface{}) error {
	var validationErrors ValidationErrors

	vType := reflect.TypeOf(v)
	vValue := reflect.ValueOf(v)

	// Валидировать будем только структуры
	if vType.Kind() != reflect.Struct {
		return nil
	}

	for i := 0; i < vType.NumField(); i++ {
		f := vType.Field(i)
		v := vValue.Field(i)
		tag := f.Tag.Get("validate")

		if tag == "" {
			continue
		}

		if !isAllowedType(f.Type) {
			continue
		}

		validators, err := parseTag(tag)
		if err != nil {
			return err
		}

		for _, validator := range validators {
			if f.Type.Kind() == reflect.Slice {
				for j := 0; j < vValue.Field(i).Len(); j++ {
					fieldFullName := fmt.Sprintf("%s.%d", f.Name, j)
					validateField(fieldFullName, vValue.Field(i).Index(j), validator, &validationErrors)
				}
			} else {
				validateField(f.Name, v, validator, &validationErrors)
			}
		}
	}

	if len(validationErrors) == 0 {
		return nil
	}

	return validationErrors
}

func validateField(name string, v reflect.Value, validator Validator, errorBucket *ValidationErrors) {
	validationErr := validator.Validate(v)
	if validationErr != nil {
		var va ValidationErrors
		if errors.As(validationErr, &va) {
			*errorBucket = append(*errorBucket, va...)
		} else {
			*errorBucket = append(*errorBucket, ValidationError{name, validationErr})
		}
	}
}

func isAllowedType(t reflect.Type) bool {
	return t.Kind() == reflect.Int ||
		t.Kind() == reflect.String ||
		t == reflect.SliceOf(reflect.TypeOf("")) || // SliceOf(string)
		t == reflect.SliceOf(reflect.TypeOf(1)) || // SliceOf(int)
		t.Kind() == reflect.Struct
}

func parseTag(tag string) ([]Validator, error) {
	var validators []Validator

	parts := strings.Split(tag, "|")
	for _, part := range parts {
		optionParts := strings.Split(part, ":")

		if len(optionParts) == 0 {
			return nil, fmt.Errorf("failed to parse validator: %s", part)
		}

		var (
			validatorName   string
			validatorOption string
		)

		if len(optionParts) > 0 {
			validatorName = optionParts[0]
		}

		if len(optionParts) > 1 {
			validatorOption = strings.TrimSpace(optionParts[1])
		}

		var validator Validator
		var err error

		switch validatorName {
		case "len":
			validator, err = NewStrLenValidator(validatorOption)
		case "min":
			validator, err = NewNumMinValidator(validatorOption)
		case "max":
			validator, err = NewNumMaxValidator(validatorOption)
		case "regexp":
			validator, err = NewRegExpValidator(validatorOption)
		case "in":
			validator, err = NewStrEnumValidator(validatorOption)
		case "nested":
			validator, err = NewNestedValidator(validatorOption)
		default:
			validator = nil
		}

		if err != nil {
			return nil, fmt.Errorf("failed to create validator %s: %w", validatorName, err)
		}

		if validator != nil {
			validators = append(validators, validator)
		}
	}

	return validators, nil
}

type Validator interface {
	Validate(v reflect.Value) error
}

type StrLenValidator struct {
	len int
}

func NewStrLenValidator(options string) (*StrLenValidator, error) {
	var (
		length int
		err    error
	)

	if length, err = strconv.Atoi(options); err != nil {
		return nil, fmt.Errorf("%s is not a number", options)
	}

	return &StrLenValidator{length}, nil
}

func (v *StrLenValidator) Validate(val reflect.Value) error {
	strVal := val.String()
	if len(strVal) != v.len {
		return ErrStrLen
	}

	return nil
}

type NumMinvalidator struct {
	min int64
}

func NewNumMinValidator(options string) (*NumMinvalidator, error) {
	var (
		num int
		err error
	)

	if num, err = strconv.Atoi(options); err != nil {
		return nil, fmt.Errorf("%s is not a number", options)
	}

	return &NumMinvalidator{int64(num)}, nil
}

func (v *NumMinvalidator) Validate(val reflect.Value) error {
	if intVal := val.Int(); intVal < v.min {
		return ErrNumRange
	}

	return nil
}

type NumMaxValidator struct {
	max int64
}

func NewNumMaxValidator(options string) (*NumMaxValidator, error) {
	var (
		num int
		err error
	)

	if num, err = strconv.Atoi(options); err != nil {
		return nil, fmt.Errorf("%s is not a number", options)
	}

	return &NumMaxValidator{int64(num)}, nil
}

func (v *NumMaxValidator) Validate(val reflect.Value) error {
	if intVal := val.Int(); intVal > v.max {
		return ErrNumRange
	}

	return nil
}

type RegExpValidator struct {
	re *regexp.Regexp
}

func NewRegExpValidator(options string) (*RegExpValidator, error) {
	re, err := regexp.Compile(options)
	if err != nil {
		return nil, fmt.Errorf("failed to parse regexp: %w", err)
	}

	return &RegExpValidator{re}, nil
}

func (v *RegExpValidator) Validate(val reflect.Value) error {
	value := val.String()

	if !v.re.Match([]byte(value)) {
		return ErrRegexp
	}

	return nil
}

type StrEnumValidator struct {
	enums []string
}

func NewStrEnumValidator(options string) (*StrEnumValidator, error) {
	return &StrEnumValidator{strings.Split(options, ",")}, nil
}

func (v *StrEnumValidator) Validate(val reflect.Value) error {
	var value string

	if val.Kind() == reflect.Int {
		value = strconv.Itoa(int(val.Int()))
	} else {
		value = val.String()
	}

	for _, expected := range v.enums {
		if value == expected {
			return nil
		}
	}

	return ErrStrEnum
}

type NestedValidator struct{}

func NewNestedValidator(options string) (*NestedValidator, error) {
	return &NestedValidator{}, nil
}

func (v *NestedValidator) Validate(val reflect.Value) error {
	return Validate(val.Interface())
}
