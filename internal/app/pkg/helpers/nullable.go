package helpers

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	// ErrUnsupportedConversion is an error that occurs when attempting to convert a value to an unsupported type.
	// This typically happens when Scan is called with a value that cannot be converted to the target type T.
	ErrUnsupportedConversion = errors.New("unsupported type conversion")
)

// Nullable is a generic struct that holds a nullable value of any type T.
// It keeps track of the value (Val), a flag (Valid) indicating whether the value has been set and a flag (Present)
// indicating if the value is in the struct.
// This allows for better handling of nullable and undefined values, ensuring proper value management and serialization.
type Nullable[T any] struct {
	Val     T
	Valid   bool
	Present bool
}

// NewNullable creates a new Nullable with the given value and sets Valid to true.
// This is useful when you want to create a Nullable with an initial value, explicitly marking it as set.
func NewNullable[T any](value T) Nullable[T] {
	return Nullable[T]{Val: value, Valid: true, Present: true}
}

// Scan implements the bson.Unmarshaler interface for Nullable, allowing it to be used as a nullable field in MongoDB operations.
// It is responsible for properly setting the Valid flag and converting the scanned value to the target type T.
// This enables seamless integration with go.mongodb.org/mongo-driver when working with nullable values.
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

// Value implements the bson.Marshaler interface for Nullable, enabling it to be used as a nullable field in MongoDB operations.
// This method ensures that the correct value is returned for serialization, handling unset Nullable values by returning nil.
func (n Nullable[T]) MarshalBSON() ([]byte, error) {
	if !n.Valid {
		return nil, nil
	}

	return bson.Marshal(n.Val)
}

// UnmarshalJSON implements the json.Unmarshaler interface for Nullable, allowing it to be used as a nullable field in JSON operations.
// This method ensures proper unmarshalling of JSON data into the Nullable value, correctly setting the Valid flag based on the JSON data.
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
// This method ensures proper marshalling of Nullable values into JSON data, representing unset values as null in the serialized output.
func (n Nullable[T]) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return []byte("null"), nil
	}

	return json.Marshal(n.Val)
}

// OrElse returns the underlying Val if valid otherwise returns the provided defaultVal
func (n Nullable[T]) OrElse(defaultVal T) T {
	if n.Valid {
		return n.Val
	} else {
		return defaultVal
	}
}

// zeroValue is a helper function that returns the zero value for the generic type T.
// It is used to set the zero value for the Val field of the Nullable struct when the value is nil.
func zeroValue[T any]() T {
	var zero T
	return zero
}

// convertToType is a helper function that attempts to convert the given value to type T.
// This function is used by Scan to properly handle value conversion, ensuring that Nullable values are always of the correct type.
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

	isNumeric := func(kind reflect.Kind) bool {
		return kind >= reflect.Int && kind <= reflect.Float64
	}

	// Check if the value is a numeric type and if T is also a numeric type.
	if isNumeric(valueType.Kind()) && isNumeric(targetType.Kind()) {
		convertedValue := reflect.ValueOf(value).Convert(targetType)
		return convertedValue.Interface().(T), nil
	}

	return zero, ErrUnsupportedConversion
}

// FindOne is a generic function that finds a single document in the specified MongoDB collection that matches the given filter.
// It returns the found document as a Nullable value of type T.
// This function incorporates security measures by using context and supporting BSON types.
func FindOne[T any](ctx context.Context, collection *mongo.Collection, filter interface{}) (Nullable[T], error) {
	var result T
	err := collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return Nullable[T]{}, nil
		}
		return Nullable[T]{}, err
	}
	return NewNullable(result), nil
}

// InsertOne is a generic function that inserts a single document of type T into the specified MongoDB collection.
// It returns the ID of the inserted document as a Nullable value.
// This function incorporates security measures by using context and supporting BSON types.
func InsertOne[T any](ctx context.Context, collection *mongo.Collection, document T) (Nullable[primitive.ObjectID], error) {
	result, err := collection.InsertOne(ctx, document)
	if err != nil {
		return Nullable[primitive.ObjectID]{}, err
	}
	return NewNullable[primitive.ObjectID](result.InsertedID.(primitive.ObjectID)), nil
}

// UpdateOne is a generic function that updates a single document in the specified MongoDB collection that matches the given filter.
// It applies the provided update to the document and returns the number of modified documents.
// This function incorporates security measures by using context and supporting BSON types.
func UpdateOne[T any](ctx context.Context, collection *mongo.Collection, filter interface{}, update interface{}) (int64, error) {
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return 0, err
	}
	return result.ModifiedCount, nil
}

// DeleteOne is a generic function that deletes a single document in the specified MongoDB collection that matches the given filter.
// It returns the number of deleted documents.
// This function incorporates security measures by using context and supporting BSON types.
func DeleteOne(ctx context.Context, collection *mongo.Collection, filter interface{}) (int64, error) {
	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return 0, err
	}
	return result.DeletedCount, nil
}
