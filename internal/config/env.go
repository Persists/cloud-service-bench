package config

import (
	"fmt"
	"os"
)

func GetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		fmt.Println(key + " environment variable is not set")
		os.Exit(1)
	}
	return value
}
