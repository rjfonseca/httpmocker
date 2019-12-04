package main

import (
	"fmt"
	"os"
	"strconv"
)

var (
	prefix = ""
)

func getEnvOrString(key string, defaultValue string) string {
	if envValue, ok := os.LookupEnv(key); ok {
		return envValue
	}
	return defaultValue
}

func getEnvOrBool(key string, defaultValue bool) bool {
	if envValue, ok := os.LookupEnv(key); ok {
		if boolValue, err := strconv.ParseBool(envValue); err != nil {
			panic(fmt.Errorf("Environment variable %s is invalid: %s", key, err))
		} else {
			return boolValue
		}
	}
	return defaultValue
}

func getEnvOrInt(key string, defaultValue int) int {
	if envValue, ok := os.LookupEnv(key); ok {
		if intValue, err := strconv.Atoi(envValue); err != nil {
			panic(fmt.Errorf("Environment variable %s is invalid: %s", key, err))
		} else {
			return intValue
		}
	}
	return defaultValue
}
