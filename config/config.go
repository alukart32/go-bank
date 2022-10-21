package config

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"gopkg.in/yaml.v3"
)

const (
	projectDirName = "bank"
	Default        = "default"
	Test           = "test"
)

type (
	DB struct {
		URI    string `env-required:"true" env:"DB_URI"`
		Driver string `env:"DB_DRIVER" env-default:"postgres"`
		DBname string `env-required:"true" env:"DB_NAME"`
	}

	HTTP struct {
		Port            string        `env:"HTTP_PORT" env-default:"80"`
		ReadTimeout     time.Duration `env:"HTTP_READ_TIMEOUT" env-default:"5000us"`
		WriteTimeout    time.Duration `env:"HTTP_WRITE_TIMEOUT" env-default:"5000us"`
		ShutdownTimeout time.Duration `env:"HTTP_SHUTDOWN_TIMEOUT" env-default:"3000us"`
	}
	Logger struct {
	}

	Config struct {
		DB     DB
		Http   HTTP
		Logger Logger
	}
)

var (
	profiles = map[string]string{
		"test":    "test.env",
		"default": ".env",
	}
)

func New(profile string) (Config, error) {
	var cfg Config

	cfgPath := fmt.Sprint(string(getRootDir()), "/config/", profiles["default"])
	if v, ok := profiles[profile]; ok {
		cfgPath = fmt.Sprint(string(getRootDir()), "/config/", v)
	}

	err := cleanenv.ReadConfig(cfgPath, &cfg)
	if err != nil {
		return cfg, fmt.Errorf("config error: %w", err)
	}

	err = cleanenv.ReadEnv(&cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}

// Print outputs the configuration in YAML format.
func Print(c Config) {
	if data, err := yaml.Marshal(&c); err != nil {
		log.Println("can not print config")
	} else {
		log.Printf("config data\n%s%v", "---\n", string(data))
	}
}

// getRootDir get project root dir.
func getRootDir() []byte {
	re := regexp.MustCompile(`^(.*` + projectDirName + `)`)
	cwd, _ := os.Getwd()
	return re.Find([]byte(cwd))
}
