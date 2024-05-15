package helpers

import (
	"encoding/json"
	"errors"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
)

var (
	// ErrUnsupportedConversion is an error that occurs when attempting to convert a value to an unsupported type.
	ErrUnsupportedConversion = errors.New("unsupported type conversion")
)

// Nullable is a generic struct that holds a nullable value of any type T.
// It keeps track of the value (Val), a flag (Valid) indicating whether the value has been set, and a flag (Present)
// indicating if the value is in the struct.
type Nullable[T any] struct {
	Val     T
	Valid   bool
	Present bool
}

// NewNullable creates a new Nullable with the given value and sets Valid and Present to true.
func NewNullable[T any](value T) Nullable[T] {
	return Nullable[T]{Val: value, Valid: true, Present: true}
}

// UnmarshalBSON implements the bson.Unmarshaler interface for Nullable, allowing it to be used as a nullable field in MongoDB operations.
func (n *Nullable[T]) UnmarshalBSON(data []byte) error {
	n.Present = true

	if len(data) == 0 {
		n.Val = zeroValue[T]()
		n.Valid = false
		return nil
	}

	var value T
	if err := bson.Unmarshal(data, &value); err != nil {
		return err
	}

	n.Val = value
	n.Valid = true
	return nil
}

// MarshalBSON implements the bson.Marshaler interface for Nullable, enabling it to be used as a nullable field in MongoDB operations.
func (n Nullable[T]) MarshalBSON() ([]byte, error) {
	if !n.Valid {
		return nil, nil
	}

	return bson.Marshal(n.Val)
}

// UnmarshalJSON implements the json.Unmarshaler interface for Nullable, allowing it to be used as a nullable field in JSON operations.
func (n *Nullable[T]) UnmarshalJSON(data []byte) error {
	n.Present = true

	if string(data) == "null" {
		n.Valid = false
		return nil
	}

	var value T
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	n.Val = value
	n.Valid = true
	return nil
}

// MarshalJSON implements the json.Marshaler interface for Nullable, enabling it to be used as a nullable field in JSON operations.
func (n Nullable[T]) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return json.Marshal(nil)
	}

	return json.Marshal(n.Val)
}

// OrElse returns the underlying Val if valid; otherwise, it returns the provided defaultVal.
func (n Nullable[T]) OrElse(defaultVal T) T {
	if n.Valid {
		return n.Val
	}
	return defaultVal
}

// zeroValue is a helper function that returns the zero value for the generic type T.
func zeroValue[T any]() T {
	var zero T
	return zero
}

// ConvertToType is a helper function that attempts to convert the given value to type T.
func ConvertToType[T any](value any) (T, error) {
	var zero T
	if value == nil {
		return zero, nil
	}

	valueType := reflect.TypeOf(value)
	targetType := reflect.TypeOf(zero)
	if valueType == targetType {
		return value.(T), nil
	}

	if isNumeric(valueType.Kind()) && isNumeric(targetType.Kind()) {
		convertedValue := reflect.ValueOf(value).Convert(targetType)
		return convertedValue.Interface().(T), nil
	}

	return zero, ErrUnsupportedConversion
}

// isNumeric is a helper function that checks if a reflect.Kind is numeric.
func isNumeric(kind reflect.Kind) bool {
	return kind >= reflect.Int && kind <= reflect.Float64
}
