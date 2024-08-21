package helpers

import (
	"encoding/json"
	"errors"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
)

var (
	ErrUnsupportedConversion = errors.New("unsupported type conversion")
	ErrNullValue             = errors.New("null value")
)

// Nullable is a generic struct that holds a nullable value of any type T.
type Nullable[T any] struct {
	Value   T
	Present bool
}

// NewNullable creates a new Nullable with the given value.
func NewNullable[T any](value T) Nullable[T] {
	v := reflect.ValueOf(value)

	// Check if the value is the zero value for its type
	if !v.IsValid() || (v.IsZero() && v.Kind() != reflect.Bool) {
		return Nullable[T]{Value: value, Present: false}
	}

	return Nullable[T]{Value: value, Present: true}
}

// Get returns the value and a boolean indicating if the value is present.
func (n Nullable[T]) Get() (T, bool) {
	return n.Value, n.Present
}

// Set sets the value and marks it as present.
func (n *Nullable[T]) Set(value T) {
	n.Value = value
	n.Present = true
}

// Clear clears the value and marks it as not present.
func (n *Nullable[T]) Clear() {
	var zero T
	n.Value = zero
	n.Present = false
}

// UnmarshalBSON implements the bson.Unmarshaler interface.
func (n *Nullable[T]) UnmarshalBSON(data []byte) error {
	if len(data) == 0 {
		n.Clear()
		return nil
	}

	var value T
	if err := bson.Unmarshal(data, &value); err != nil {
		return err
	}

	n.Set(value)
	return nil
}

// MarshalBSON implements the bson.Marshaler interface.
func (n Nullable[T]) MarshalBSON() ([]byte, error) {
	if !n.Present {
		return bson.Marshal(nil)
	}
	return bson.Marshal(n.Value)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (n *Nullable[T]) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		n.Clear()
		return nil
	}

	var value T
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	n.Set(value)
	return nil
}

// MarshalJSON implements the json.Marshaler interface.
func (n Nullable[T]) MarshalJSON() ([]byte, error) {
	if !n.Present {
		return json.Marshal(nil)
	}
	return json.Marshal(n.Value)
}

// OrElse returns the value if present, otherwise returns the provided default value.
func (n Nullable[T]) OrElse(defaultVal T) T {
	if n.Present {
		return n.Value
	}
	return defaultVal
}

// Map applies the provided function to the value if present.
func (n Nullable[T]) Map(f func(T) T) Nullable[T] {
	if !n.Present {
		return n
	}
	return NewNullable(f(n.Value))
}

// FlatMap applies the provided function to the value if present, returning a new Nullable.
func (n Nullable[T]) FlatMap(f func(T) Nullable[T]) Nullable[T] {
	if !n.Present {
		return n
	}
	return f(n.Value)
}

// Filter returns the Nullable if the predicate is true, otherwise returns an empty Nullable.
func (n Nullable[T]) Filter(predicate func(T) bool) Nullable[T] {
	if !n.Present || !predicate(n.Value) {
		return Nullable[T]{}
	}
	return n
}

// ConvertToType is a helper function that attempts to convert the given value to type T.
func ConvertToType[T any](value interface{}) (T, error) {
	var zero T
	if value == nil {
		return zero, ErrNullValue
	}

	valueType := reflect.TypeOf(value)
	targetType := reflect.TypeOf(zero)

	if valueType.ConvertibleTo(targetType) {
		convertedValue := reflect.ValueOf(value).Convert(targetType)
		return convertedValue.Interface().(T), nil
	}

	return zero, ErrUnsupportedConversion
}

// IsNumeric is a helper function that checks if a reflect.Kind is numeric.
func IsNumeric(kind reflect.Kind) bool {
	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return true
	default:
		return false
	}
}
