package function

func Coalesce(value, defaultValue interface{}) interface{} {
	if value == nil || value == "" || value == 0 || value == 0.0 {
		return defaultValue
	}
	return value
}
