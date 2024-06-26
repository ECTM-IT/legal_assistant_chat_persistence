package helpers

import (
	"encoding/json"
	"errors"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
)

var (
	ErrUnsupportedConversion = errors.New("unsupported type conversion")
	ErrNullValue             = errors.New("null Value")
)

// Nullable is a generic struct that holds a nullable Value of any type T.
type Nullable[T any] struct {
	Value   T
	Present bool
}

// NewNullable creates a new Nullable with the given Value.
func NewNullable[T any](Value T) Nullable[T] {
	return Nullable[T]{Value: Value, Present: true}
}

// Get returns the Value and a boolean indicating if the Value is Present.
func (n Nullable[T]) Get() (T, bool) {
	return n.Value, n.Present
}

// Set sets the Value and marks it as Present.
func (n *Nullable[T]) Set(Value T) {
	n.Value = Value
	n.Present = true
}

// Clear clears the Value and marks it as not Present.
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

	var Value T
	if err := bson.Unmarshal(data, &Value); err != nil {
		return err
	}

	n.Set(Value)
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

	var Value T
	if err := json.Unmarshal(data, &Value); err != nil {
		return err
	}

	n.Set(Value)
	return nil
}

// MarshalJSON implements the json.Marshaler interface.
func (n Nullable[T]) MarshalJSON() ([]byte, error) {
	if !n.Present {
		return json.Marshal(nil)
	}
	return json.Marshal(n.Value)
}

// OrElse returns the Value if Present, otherwise returns the provided default Value.
func (n Nullable[T]) OrElse(defaultVal T) T {
	if n.Present {
		return n.Value
	}
	return defaultVal
}

// Map applies the provided function to the Value if Present.
func (n Nullable[T]) Map(f func(T) T) Nullable[T] {
	if !n.Present {
		return n
	}
	return NewNullable(f(n.Value))
}

// FlatMap applies the provided function to the Value if Present, returning a new Nullable.
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

// ConvertToType is a helper function that attempts to convert the given Value to type T.
func ConvertToType[T any](Value any) (T, error) {
	var zero T
	if Value == nil {
		return zero, ErrNullValue
	}

	ValueType := reflect.TypeOf(Value)
	targetType := reflect.TypeOf(zero)
	if ValueType == targetType {
		return Value.(T), nil
	}

	if isNumeric(ValueType.Kind()) && isNumeric(targetType.Kind()) {
		convertedValue := reflect.ValueOf(Value).Convert(targetType)
		return convertedValue.Interface().(T), nil
	}

	return zero, ErrUnsupportedConversion
}

// isNumeric is a helper function that checks if a reflect.Kind is numeric.
func isNumeric(kind reflect.Kind) bool {
	return kind >= reflect.Int && kind <= reflect.Float64
}
