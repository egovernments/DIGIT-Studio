package config

import (
	"os"
)

/*
func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}*/

func GetEnv(key string) string {
	return os.Getenv(key)
}
