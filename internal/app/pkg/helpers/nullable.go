package helpers

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
)

var (
	ErrUnsupportedConversion = errors.New("unsupported type conversion")
	ErrNullValue             = errors.New("null value")
)

// Nullable represents a generic nullable type for any type T.
// It encapsulates a value and a flag indicating its presence.
type Nullable[T any] struct {
	Value   T    `json:"value,omitempty" bson:"value,omitempty"`
	Present bool `json:"present" bson:"present"`
}

// NewNullable creates a new Nullable instance with the provided value.
// It determines the presence based on whether the value is the zero value for its type.
func NewNullable[T any](value T) Nullable[T] {
	var zero T
	isPresent := !isZeroValue(value, zero)
	return Nullable[T]{Value: value, Present: isPresent}
}

// Get retrieves the value and a boolean indicating its presence.
func (n Nullable[T]) Get() (T, bool) {
	return n.Value, n.Present
}

// Set assigns a new value and marks it as present.
func (n *Nullable[T]) Set(value T) {
	n.Value = value
	n.Present = true
}

// Clear resets the value to its zero value and marks it as not present.
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
		return bson.Marshal(bson.M{"present": false})
	}
	return bson.Marshal(bson.M{"value": n.Value, "present": n.Present})
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

// OrElse returns the value if present; otherwise, it returns the provided default value.
func (n Nullable[T]) OrElse(defaultVal T) T {
	if n.Present {
		return n.Value
	}
	return defaultVal
}

// Map applies the provided function to the value if present and returns a new Nullable.
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

// Filter returns the Nullable if the predicate is true; otherwise, it returns an empty Nullable.
func (n Nullable[T]) Filter(predicate func(T) bool) Nullable[T] {
	if !n.Present || !predicate(n.Value) {
		return Nullable[T]{}
	}
	return n
}

// ConvertToType attempts to convert the given value to type T.
// It returns an error if the conversion is not supported.
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

// IsNumeric checks if a reflect.Kind is numeric.
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

// String provides a string representation of the Nullable.
func (n Nullable[T]) String() string {
	if !n.Present {
		return "null"
	}
	return fmt.Sprintf("%v", n.Value)
}

// isZeroValue checks whether the provided value is the zero value for its type.
// This helper function improves readability and reuse.
func isZeroValue[T any](value, zero T) bool {
	return reflect.DeepEqual(value, zero)
}
