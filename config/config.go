package config

import (
	"os"
	"regexp"

	"github.com/joho/godotenv"
)

type database struct {
	URL string
}

type Config struct {
	Database database
}

func LoadEnv(filename string) {
	re := regexp.MustCompile(`^(.*` + "sweeper" + `)`)

	cwd, _ := os.Getwd()
	rootPath := re.Find([]byte(cwd))

	err := godotenv.Load(string(rootPath) + `/` + filename)
	if err != nil {
		godotenv.Load()
	}
}

func New() *Config {

	return &Config{
		Database: database{
			URL: os.Getenv("DATABASE_URL"),
		},
	}
}
