package util

import "os"

func GetEnv(key string, deflt string) string {
	value := os.Getenv(key)
	if value == "" {
		return deflt
	}

	return value
}
