package utils

import (
	"fmt"
)

// ToString - Basic type conversion functions
func ToString(val any) (string, bool) {
	switch v := val.(type) {
	case string:
		return v, true
	case []byte:
		return string(v), true
	case fmt.Stringer:
		return v.String(), true
	default:
		return fmt.Sprintf("%v", val), true
	}
}
