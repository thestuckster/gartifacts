package internal

import (
	"github.com/joho/godotenv"
	"os"
)

func LoadEnvironment() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
}

func GetEnvVar(key string) string {
	env := os.Getenv(key)
	if env == "" {
		panic("Environment variable " + key + " not set")
	}

	return env
}
