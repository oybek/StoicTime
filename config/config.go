package config

import (
	"errors"
	"os"
)

type Config struct {
	OpenAIToken string
}

func NewConfig() (c Config, err error) {
	c.OpenAIToken, err = getenv("OPENAI_TOKEN")
	if err != nil {
		return
	}

	return c, nil
}

func getenv(key string) (string, error) {
	value, exists := os.LookupEnv(key)
	if !exists {
		return "", errors.New(key + " env variable is not set")
	}
	return value, nil
}
