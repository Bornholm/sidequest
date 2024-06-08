package env

import (
	"os"
	"strconv"
)

func Bool(key string, defaultValue bool) bool {
	raw, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}

	value, err := strconv.ParseBool(raw)
	if err != nil {
		return defaultValue
	}

	return value
}

func Int(key string, defaultValue int) int {
	raw, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}

	value, err := strconv.ParseInt(raw, 10, 32)
	if err != nil {
		return defaultValue
	}

	return int(value)
}

func String(key string, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}

	return value
}
