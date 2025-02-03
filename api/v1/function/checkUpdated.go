package function

import "reflect"

func Coalesce(value, defaultValue interface{}) interface{} {
	if value == nil {
		return defaultValue
	}

	v := reflect.ValueOf(value)
	// Handle string and string-based types
	if v.Kind() == reflect.String && v.String() == "" {
		return defaultValue
	}
	// Handle numeric types
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if v.Int() == 0 {
			return defaultValue
		}
	case reflect.Float32, reflect.Float64:
		if v.Float() == 0 {
			return defaultValue
		}
	}
	return value
}
