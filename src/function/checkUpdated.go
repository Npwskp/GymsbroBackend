package function

func Coalesce(value, defaultValue interface{}) interface{} {
	if value == nil {
		return defaultValue
	}
	return value
}
