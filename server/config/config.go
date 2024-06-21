package config

import (
	"errors"
	"os"

	"github.com/abibby/salusa/env"
	"github.com/joho/godotenv"
)

type Config struct {
	Port     int
	BasePath string

	HIDPort   int
	AdapterID string
}

func Load() *Config {
	err := godotenv.Load("./.env")
	if errors.Is(err, os.ErrNotExist) {
		// fall through
	} else if err != nil {
		panic(err)
	}

	return &Config{
		Port:      env.Int("PORT", 3808),
		HIDPort:   env.Int("HID_PORT", 38808),
		AdapterID: env.String("BLUETOOTH_ADAPTER_ID", "hci0"),
		BasePath:  env.String("BASE_PATH", ""),
	}
}

func (c *Config) GetHTTPPort() int {
	return c.Port
}
func (c *Config) GetBaseURL() string {
	return c.BasePath
}
