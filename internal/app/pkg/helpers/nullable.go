package helpers

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
)

// Predefined errors for clarity and consistency.
var (
	ErrUnsupportedConversion = errors.New("unsupported type conversion")
	ErrNullValue             = errors.New("null value provided")
)

// Nullable represents a generic optional (nullable) value of type T.
// It holds a value and a boolean flag indicating whether the value is present.
//
// The Nullable type is useful when dealing with optional fields, especially
// in scenarios involving database queries, JSON/BSON serialization, and
// data transformations where a value may or may not be present.
type Nullable[T any] struct {
	// Value holds the actual data. If Present is false, this value should be
	// considered invalid or ignored.
	Value T `json:"value,omitempty" bson:"value,omitempty"`

	// Present indicates whether the value is actually set.
	Present bool `json:"present" bson:"present"`
}

// NewNullable creates a new Nullable instance with the provided value.
// The presence flag is determined by checking whether the value is its
// zero-value. If not zero, Present is true; otherwise false.
//
// Example:
//
//	n := NewNullable(42)   // Present = true
//	m := NewNullable(int(0)) // Present = false
func NewNullable[T any](value T) Nullable[T] {
	isPresent := !isZeroValue(value)
	return Nullable[T]{Value: value, Present: isPresent}
}

// Get returns the underlying value and a boolean indicating its presence.
func (n Nullable[T]) Get() (T, bool) {
	return n.Value, n.Present
}

// Set assigns a new value and marks the Nullable as present.
func (n *Nullable[T]) Set(value T) {
	n.Value = value
	n.Present = true
}

// Clear resets the Nullable to its zero value and marks it as not present.
func (n *Nullable[T]) Clear() {
	var zero T
	n.Value = zero
	n.Present = false
}

// UnmarshalBSON implements the bson.Unmarshaler interface, allowing the Nullable
// to be correctly populated from BSON data.
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

// MarshalBSON implements the bson.Marshaler interface, allowing the Nullable
// to be correctly serialized to BSON. If not present, it sets 'present' to false.
func (n Nullable[T]) MarshalBSON() ([]byte, error) {
	if !n.Present {
		// Represent absence in BSON as {present: false}
		return bson.Marshal(bson.M{"present": false})
	}
	return bson.Marshal(bson.M{"value": n.Value, "present": n.Present})
}

// UnmarshalJSON implements the json.Unmarshaler interface. If the input is "null",
// the Nullable is cleared. Otherwise, the value is unmarshaled and marked as present.
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

// MarshalJSON implements the json.Marshaler interface. If not present,
// null is returned. Otherwise, the value is marshaled normally.
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
// If the value is not present, it returns a new empty Nullable.
func (n Nullable[T]) Map(f func(T) T) Nullable[T] {
	if !n.Present {
		return Nullable[T]{}
	}
	return NewNullable(f(n.Value))
}

// FlatMap applies the provided function to the value if present, returning a new Nullable.
// If the value is not present, it returns a new empty Nullable.
func (n Nullable[T]) FlatMap(f func(T) Nullable[T]) Nullable[T] {
	if !n.Present {
		return Nullable[T]{}
	}
	return f(n.Value)
}

// Filter returns the Nullable if the predicate is true; otherwise, it returns an empty Nullable.
// This allows for conditionally retaining values.
func (n Nullable[T]) Filter(predicate func(T) bool) Nullable[T] {
	if !n.Present || !predicate(n.Value) {
		return Nullable[T]{}
	}
	return n
}

// ConvertToType attempts to convert the given interface{} value to the specified generic type T.
// If the conversion is not possible or the value is nil, it returns an error.
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

// IsNumeric checks if a reflect.Kind represents a numeric type.
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

// String returns a string representation of the Nullable.
// If the value is not present, "null" is returned.
func (n Nullable[T]) String() string {
	if !n.Present {
		return "null"
	}
	return fmt.Sprintf("%v", n.Value)
}

// isZeroValue checks whether the provided value is the zero value for its type.
// Using reflect.ValueOf(value).IsZero() is preferred over reflect.DeepEqual
// for performance and clarity. This requires Go 1.13+.
func isZeroValue[T any](value T) bool {
	return reflect.ValueOf(value).IsZero()
}
