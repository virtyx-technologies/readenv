// Package readenv provides an easy interface for reading environment variables
// into an options struct by adding tags to the struct fields.
//
// To use readenv, simply add tags to your struct containing the environment
// variable the field should be read from. For example:
//
//     type options struct {
//       Port int `env:"PORT"`
//     }
//
// Note that the field must be exported so that it can be writeable.
package readenv

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

// ReadEnv reads environment variables into the provided struct pointer. If
// there are any problems reading or parsing the environment variables, an error
// is returned.
//
// The argument to ReadEnv must be a pointer to a struct. If any other value is
// passed, an error will be returned.
//
// Some special values are recognized when parsing environment variables into
// boolean fields: if the environment variable is set to "no", "off", "0", or is
// empty, the bool will be set to false. Any other value in the environment will
// set it to true.
func ReadEnv(dest interface{}) error {
	v := reflect.ValueOf(dest)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("readenv: dest should be pointer to struct, but was %v", v.Type())
	}
	v = v.Elem()
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)
		if err := readField(fieldValue, field); err != nil {
			return fmt.Errorf("readenv: could not set %s: %v", field.Name, err)
		}
	}
	return nil
}

func readField(val reflect.Value, field reflect.StructField) error {
	if !val.CanSet() {
		return fmt.Errorf("field is not writeable")
	}
	if envName, ok := field.Tag.Lookup("env"); ok {
		if isInt(field.Type) {
			if err := readEnvInt(val, envName); err != nil {
				return err
			}
		} else if isString(field.Type) {
			if err := readEnvString(val, envName); err != nil {
				return err
			}
		} else if isFloat(field.Type) {
			if err := readEnvFloat(val, envName); err != nil {
				return err
			}
		} else if isBool(field.Type) {
			readEnvBool(val, envName)
		}
	}
	return nil
}

func isString(t reflect.Type) bool {
	return t == reflect.TypeOf("")
}

func readEnvString(field reflect.Value, name string) error {
	value := strings.TrimSpace(os.Getenv(name))
	if value == "" {
		return fmt.Errorf("%s is not set", name)
	}
	field.SetString(value)
	return nil
}

func isInt(t reflect.Type) bool {
	return (t == reflect.TypeOf(int(0)) ||
		t == reflect.TypeOf(int8(0)) ||
		t == reflect.TypeOf(int16(0)) ||
		t == reflect.TypeOf(int32(0)) ||
		t == reflect.TypeOf(int64(0)))
}

func readEnvInt(field reflect.Value, name string) error {
	value, err := strconv.Atoi(os.Getenv(name))
	if err != nil {
		return fmt.Errorf("%s is not a number: %v", name, err)
	}
	field.SetInt(int64(value))
	return nil
}

func isFloat(t reflect.Type) bool {
	return (t == reflect.TypeOf(float32(0)) ||
		t == reflect.TypeOf(float64(0)))
}

func readEnvFloat(field reflect.Value, name string) error {
	value, err := strconv.ParseFloat(os.Getenv(name), 64)
	if err != nil {
		return fmt.Errorf("%s is not a float: %v", name, err)
	}
	field.SetFloat(value)
	return nil
}

func isBool(t reflect.Type) bool {
	return t == reflect.TypeOf(false)
}

func readEnvBool(field reflect.Value, name string) {
	v := strings.ToLower(strings.TrimSpace(os.Getenv(name)))
	if v == "" || v == "no" || v == "off" || v == "0" {
		field.SetBool(false)
	} else {
		field.SetBool(true)
	}
}
