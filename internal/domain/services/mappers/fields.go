package mappers

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/app/pkg/helpers"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// FieldError represents an error that occurred during field assignment
type FieldError struct {
	FieldName  string
	InputValue interface{}
	ErrorType  string
	Message    string
}

func (e FieldError) Error() string {
	return fmt.Sprintf("Field: %s, Value: %v, Error: %s - %s", e.FieldName, e.InputValue, e.ErrorType, e.Message)
}

// FieldAssigner handles the assignment of fields
type FieldAssigner struct {
	ErrorHandler  func(error)
	DefaultValues map[reflect.Type]interface{}
}

// NewFieldAssigner creates a new FieldAssigner with default error handler and default values
func NewFieldAssigner() *FieldAssigner {
	return &FieldAssigner{
		ErrorHandler: func(err error) {
			fmt.Printf("Field assignment error: %v\n", err)
		},
		DefaultValues: map[reflect.Type]interface{}{
			reflect.TypeOf(""):                                     "",
			reflect.TypeOf(0):                                      0,
			reflect.TypeOf(0.0):                                    0.0,
			reflect.TypeOf(false):                                  false,
			reflect.TypeOf(time.Time{}):                            time.Time{},
			reflect.TypeOf([]interface{}{}):                        []interface{}{},
			reflect.TypeOf(map[string]interface{}{}):               map[string]interface{}{},
			reflect.TypeOf(primitive.ObjectID{}):                   primitive.NilObjectID,
			reflect.TypeOf([]primitive.ObjectID{}):                 []primitive.ObjectID{},
			reflect.TypeOf([]string{}):                             []string{},
			reflect.TypeOf([]int{}):                                []int{},
			reflect.TypeOf([]float64{}):                            []float64{},
			reflect.TypeOf(helpers.Nullable[string]{}):             helpers.Nullable[string]{Present: false},
			reflect.TypeOf(helpers.Nullable[int]{}):                helpers.Nullable[int]{Present: false},
			reflect.TypeOf(helpers.Nullable[float64]{}):            helpers.Nullable[float64]{Present: false},
			reflect.TypeOf(helpers.Nullable[bool]{}):               helpers.Nullable[bool]{Present: false},
			reflect.TypeOf(helpers.Nullable[time.Time]{}):          helpers.Nullable[time.Time]{Present: false},
			reflect.TypeOf(helpers.Nullable[primitive.ObjectID]{}): helpers.Nullable[primitive.ObjectID]{Present: false},
		},
	}
}

// AssignField assigns a value to a field, performing necessary checks and conversions
func (fa *FieldAssigner) AssignField(field reflect.Value, value interface{}, fieldName string) {
	defer func() {
		if r := recover(); r != nil {
			fa.ErrorHandler(FieldError{
				FieldName:  fieldName,
				InputValue: value,
				ErrorType:  "Panic",
				Message:    fmt.Sprintf("Recovered from panic: %v", r),
			})
		}
	}()

	if !field.CanSet() {
		fa.ErrorHandler(FieldError{
			FieldName:  fieldName,
			InputValue: value,
			ErrorType:  "Unsettable",
			Message:    "Field cannot be set",
		})
		return
	}

	if value == nil {
		fa.setDefaultValue(field, fieldName)
		return
	}

	valueType := reflect.TypeOf(value)
	fieldType := field.Type()

	// Handle Nullable types
	if fieldType.Name() == "Nullable" {
		fa.handleNullableField(field, value, fieldName)
		return
	}

	// Handle slice types
	if fieldType.Kind() == reflect.Slice {
		fa.handleSliceField(field, value, fieldName)
		return
	}

	// Type checking and coercion
	if valueType != fieldType {
		coercedValue, err := fa.coerceType(value, fieldType)
		if err != nil {
			fa.ErrorHandler(FieldError{
				FieldName:  fieldName,
				InputValue: value,
				ErrorType:  "TypeMismatch",
				Message:    err.Error(),
			})
			fa.setDefaultValue(field, fieldName)
			return
		}
		value = coercedValue
	}

	// Validation
	if err := fa.validateField(value, fieldName); err != nil {
		fa.ErrorHandler(err)
		fa.setDefaultValue(field, fieldName)
		return
	}

	field.Set(reflect.ValueOf(value))
}

func (fa *FieldAssigner) setDefaultValue(field reflect.Value, fieldName string) {
	defaultValue, ok := fa.DefaultValues[field.Type()]
	if !ok {
		fa.ErrorHandler(FieldError{
			FieldName: fieldName,
			ErrorType: "NoDefaultValue",
			Message:   fmt.Sprintf("No default value for type %v", field.Type()),
		})
		return
	}
	field.Set(reflect.ValueOf(defaultValue))
}

func (fa *FieldAssigner) handleNullableField(field reflect.Value, value interface{}, fieldName string) {
	nullableValue := reflect.New(field.Type()).Elem()
	valueField := nullableValue.FieldByName("Value")
	presentField := nullableValue.FieldByName("Present")
	validField := nullableValue.FieldByName("Valid")

	if value == nil {
		presentField.SetBool(false)
		validField.SetBool(false)
	} else {
		presentField.SetBool(true)
		validField.SetBool(true)
		fa.AssignField(valueField, value, fieldName)
	}

	field.Set(nullableValue)
}

func (fa *FieldAssigner) handleSliceField(field reflect.Value, value interface{}, fieldName string) {
	valueSlice, ok := value.([]interface{})
	if !ok {
		fa.ErrorHandler(FieldError{
			FieldName:  fieldName,
			InputValue: value,
			ErrorType:  "TypeMismatch",
			Message:    "Expected slice type",
		})
		fa.setDefaultValue(field, fieldName)
		return
	}

	sliceType := field.Type()
	slice := reflect.MakeSlice(sliceType, len(valueSlice), len(valueSlice))

	for i, v := range valueSlice {
		fa.AssignField(slice.Index(i), v, fmt.Sprintf("%s[%d]", fieldName, i))
	}

	field.Set(slice)
}

func (fa *FieldAssigner) coerceType(value interface{}, targetType reflect.Type) (interface{}, error) {
	switch targetType {
	case reflect.TypeOf(primitive.ObjectID{}):
		switch v := value.(type) {
		case string:
			return primitive.ObjectIDFromHex(v)
		}
	}

	switch targetType.Kind() {
	case reflect.String:
		return fmt.Sprintf("%v", value), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch v := value.(type) {
		case string:
			return strconv.ParseInt(v, 10, 64)
		case float64:
			return int64(v), nil
		}
	case reflect.Float32, reflect.Float64:
		switch v := value.(type) {
		case string:
			return strconv.ParseFloat(v, 64)
		case int:
			return float64(v), nil
		}
	case reflect.Bool:
		switch v := value.(type) {
		case string:
			return strconv.ParseBool(strings.ToLower(v))
		}
	case reflect.Struct:
		if targetType == reflect.TypeOf(time.Time{}) {
			switch v := value.(type) {
			case string:
				return time.Parse(time.RFC3339, v)
			case float64:
				return time.Unix(int64(v), 0), nil
			}
		}
	}
	return nil, fmt.Errorf("cannot coerce %T to %v", value, targetType)
}

func (fa *FieldAssigner) validateField(value interface{}, fieldName string) error {
	// Implement your validation logic here
	// This is a placeholder for demonstration
	switch v := value.(type) {
	case string:
		if len(v) == 0 {
			return FieldError{
				FieldName:  fieldName,
				InputValue: v,
				ErrorType:  "Validation",
				Message:    "String cannot be empty",
			}
		}
	case int, int8, int16, int32, int64, float32, float64:
		if reflect.ValueOf(v).Float() < 0 {
			return FieldError{
				FieldName:  fieldName,
				InputValue: v,
				ErrorType:  "Validation",
				Message:    "Number cannot be negative",
			}
		}
	case primitive.ObjectID:
		if v == primitive.NilObjectID {
			return FieldError{
				FieldName:  fieldName,
				InputValue: v,
				ErrorType:  "Validation",
				Message:    "ObjectID cannot be nil",
			}
		}
	}
	return nil
}

// AssignFields assigns values from a map to struct fields
func (fa *FieldAssigner) AssignFields(dest interface{}, src map[string]interface{}) {
	destValue := reflect.ValueOf(dest).Elem()
	destType := destValue.Type()

	for i := 0; i < destValue.NumField(); i++ {
		field := destValue.Field(i)
		fieldType := destType.Field(i)
		fieldName := fieldType.Name
		jsonTag := fieldType.Tag.Get("json")
		if jsonTag != "" {
			fieldName = strings.Split(jsonTag, ",")[0]
		}

		if value, ok := src[fieldName]; ok {
			fa.AssignField(field, value, fieldName)
		} else {
			fa.setDefaultValue(field, fieldName)
		}
	}
}
