package config

import (
	"errors"
	"os"

	"github.com/abibby/salusa/database/dialects"
	"github.com/abibby/salusa/database/dialects/sqlite"
	"github.com/abibby/salusa/env"
	"github.com/abibby/salusa/event"
	"github.com/joho/godotenv"
)

type Config struct {
	Port     int
	BasePath string

	Database dialects.Config
	Queue    event.Config

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
		Port:      env.Int("PORT", 2303),
		HIDPort:   env.Int("HID_PORT", 38808),
		AdapterID: env.String("BLUETOOTH_ADAPTER_ID", "hci0"),
		BasePath:  env.String("BASE_PATH", ""),
		Database:  sqlite.NewConfig(env.String("DATABASE_PATH", "./db.sqlite")),
		Queue:     event.NewChannelQueueConfig(),
	}
}

func (c *Config) GetHTTPPort() int {
	return c.Port
}
func (c *Config) GetBaseURL() string {
	return c.BasePath
}

func (c *Config) DBConfig() dialects.Config {
	return c.Database
}
func (c *Config) QueueConfig() event.Config {
	return c.Queue
}
