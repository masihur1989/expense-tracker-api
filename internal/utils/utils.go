package utils

import (
	"log"
	"os"
	"strconv"
)

// MustGet will return the env or panic if it is not present
func MustGet(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Panicf("ENV missing, key: %s" + k)
	}
	return v
}

// MustGetBool will return the env as boolean or panic if it is not present
func MustGetBool(k string) bool {
	v := os.Getenv(k)
	if v == "" {
		log.Panicf("ENV missing, key: %s" + k)
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		log.Panicln("ENV err: [" + k + "]\n" + err.Error())
	}
	return b
}

// IsProduction check if the current env is production
func IsProduction() bool {
	env := MustGet("SERVER_MODE")
	if env == "production" {
		return true
	}
	return false
}
