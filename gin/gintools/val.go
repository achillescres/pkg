package gintools

import (
	"fmt"
)

func ValueToString(key string, v any, ok bool) (string, error) {
	if !ok {
		return "", fmt.Errorf("error key %s is missing", key)
	}
	s, ok := v.(string)
	if !ok {
		return "", fmt.Errorf("error couldn't assert %s to string", key)
	}

	return s, nil
}
