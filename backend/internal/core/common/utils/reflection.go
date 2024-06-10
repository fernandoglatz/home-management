package utils

import (
	"encoding/json"
	"errors"
	"reflect"
)

func Instance[T any]() T {
	var object T
	json.Unmarshal([]byte("{}"), &object)
	return object
}

func GetTypeName(object any) string {
	return reflect.TypeOf(object).String()
}

func CopyStructFields(source any, destination any) error {
	sourceValue := reflect.ValueOf(source)
	destinationValue := reflect.ValueOf(destination)

	if sourceValue.IsNil() {
		return errors.New("source is nil")
	}

	if destinationValue.IsNil() {
		return errors.New("destination is nil")
	}

	if sourceValue.Kind() == reflect.Ptr {
		sourceValue = sourceValue.Elem()
	}

	if destinationValue.Kind() == reflect.Ptr {
		destinationValue = destinationValue.Elem()
	}

	for i := 0; i < sourceValue.NumField(); i++ {
		sourceField := sourceValue.Field(i)

		fieldName := sourceValue.Type().Field(i).Name
		destinationField := destinationValue.FieldByName(fieldName)

		if destinationField.IsValid() && destinationField.CanSet() {
			destinationField.Set(sourceField)
		}
	}

	return nil
}
