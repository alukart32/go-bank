package config

import (
	"log"
	"os"
	"regexp"
	"time"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

const projectDirName = "bank"

type (
	DB struct {
		URI         string `json:"uri" yaml:"uri"`
		Driver      string
		DBname      string `json:"dbname" yaml:"dbname"`
		InitTimeout time.Duration
	}

	Logger struct {
	}

	Config struct {
		DB     DB
		Logger Logger
	}
)

func New(profile string) (*Config, error) {
	if len(profile) == 0 {
		profile = getEnvPath() + "/config/.env"
	} else {
		profile = getEnvPath() + "/config/" + profile + ".env"
	}

	err := godotenv.Load(profile)
	if err != nil {
		return nil, err
	}

	db := &DB{
		URI:         os.Getenv("PGDB_URI"),
		Driver:      os.Getenv("DB_DRIVER"),
		DBname:      os.Getenv("PGDB_NAME"),
		InitTimeout: 25 * time.Microsecond,
	}

	return &Config{
		DB: *db,
	}, nil
}

// GetEnvPath define env vars from .env
func getEnvPath() string {
	re := regexp.MustCompile(`^(.*` + projectDirName + `)`)
	cwd, _ := os.Getwd()
	rootPath := re.Find([]byte(cwd))
	return string(rootPath)
}

// Print outputs the configuration in YAML format.
func Print(c *Config) {
	if data, err := yaml.Marshal(*c); err != nil {
		log.Println("can not print config")
	} else {
		log.Printf("config data\n%s%v", "---\n", string(data))
	}
}
