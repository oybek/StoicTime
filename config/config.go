package config

import (
	"errors"
	"os"
)

type Config struct {
	BotToken string
	PGURL    string
}

func NewConfig() (c Config, err error) {
	if c.BotToken, err = getenv("BOT_TOKEN"); err != nil {
		return
	}
	if c.PGURL, err = getenv("PG_URL"); err != nil {
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
