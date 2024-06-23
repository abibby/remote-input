package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/kardianos/service"
)

var Host string
var Port int

func Init(l service.Logger) {
	// exe, err := os.Executable()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// absExe, err := filepath.Abs(exe)
	// dir := filepath.Dir(absExe)
	err := godotenv.Load("C:/Program Files/RemoteInput/.env")
	if err == os.ErrNotExist {
		// do nothing
	} else if err != nil {
		log.Print(err)
	}
	Host = env("REMOTE_INPUT_HOST", "localhost")
	Port = envInt("REMOTE_INPUT_PORT", 38808)
}

func env(key, defaultValue string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}
	return val
}
func envInt(key string, defaultValue int) int {
	strVal, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}
	val, err := strconv.Atoi(strVal)
	if err != nil {
		return defaultValue
	}
	return val
}
