package expression

import (
	"fmt"
	"strings"
)

// joinFunc joins array elements with a separator.
func joinFunc(value interface{}, args []string) (interface{}, error) {
	separator := ","
	if len(args) > 0 {
		separator = args[0]
	}

	switch v := value.(type) {
	case []interface{}:
		strs := make([]string, len(v))
		for i, item := range v {
			strs[i] = fmt.Sprintf("%v", item)
		}
		return strings.Join(strs, separator), nil

	case []string:
		return strings.Join(v, separator), nil

	default:
		return fmt.Sprintf("%v", value), nil
	}
}

// firstFunc returns the first element of an array.
func firstFunc(value interface{}, args []string) (interface{}, error) {
	switch v := value.(type) {
	case []interface{}:
		if len(v) == 0 {
			return nil, nil
		}
		return v[0], nil

	case []string:
		if len(v) == 0 {
			return "", nil
		}
		return v[0], nil

	default:
		return value, nil
	}
}

// lastFunc returns the last element of an array.
func lastFunc(value interface{}, args []string) (interface{}, error) {
	switch v := value.(type) {
	case []interface{}:
		if len(v) == 0 {
			return nil, nil
		}
		return v[len(v)-1], nil

	case []string:
		if len(v) == 0 {
			return "", nil
		}
		return v[len(v)-1], nil

	default:
		return value, nil
	}
}

// lengthFunc returns the length of an array or string.
func lengthFunc(value interface{}, args []string) (interface{}, error) {
	switch v := value.(type) {
	case []interface{}:
		return len(v), nil

	case []string:
		return len(v), nil

	case string:
		return len(v), nil

	default:
		return 0, nil
	}
}

// defaultFunc returns the value or a default if nil/empty.
func defaultFunc(value interface{}, args []string) (interface{}, error) {
	if len(args) == 0 {
		return value, nil
	}

	defaultVal := args[0]

	switch v := value.(type) {
	case nil:
		return defaultVal, nil

	case string:
		if v == "" {
			return defaultVal, nil
		}
		return v, nil

	case []interface{}:
		if len(v) == 0 {
			return defaultVal, nil
		}
		return v, nil

	default:
		return value, nil
	}
}

// upperFunc converts a string to uppercase.
func upperFunc(value interface{}, args []string) (interface{}, error) {
	return strings.ToUpper(fmt.Sprintf("%v", value)), nil
}

// lowerFunc converts a string to lowercase.
func lowerFunc(value interface{}, args []string) (interface{}, error) {
	return strings.ToLower(fmt.Sprintf("%v", value)), nil
}

// trimFunc trims whitespace from a string.
func trimFunc(value interface{}, args []string) (interface{}, error) {
	return strings.TrimSpace(fmt.Sprintf("%v", value)), nil
}
