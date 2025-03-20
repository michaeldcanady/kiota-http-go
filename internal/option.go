package internal

import (
	"fmt"
	"reflect"
)

// Option is a function that can be used to set defaulting options.
type Option[T any] func(T) error

// ApplyOptions applies the given options to the given instance.
func ApplyOptions[T any](value T, options ...Option[T]) error {
	//check if value is a pointer
	if reflect.ValueOf(value).Kind() != reflect.Ptr {
		return fmt.Errorf("value must be a pointer")
	}

	// check if value is nil
	if reflect.ValueOf(value).IsNil() {
		return fmt.Errorf("value is nil")
	}

	// early return if no options are provided
	if len(options) == 0 {
		return nil
	}

	for _, option := range options {
		if err := option(value); err != nil {
			return err
		}
	}
	return nil
}
