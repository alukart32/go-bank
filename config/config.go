// Package config implements utility routines for manipulating env config file.
//
// The config package should only be used to set the general configuration of the application.
// Env params hold in the .env file. Declare this file in the root directory to configure
// the launch of the application. Meanwhile, for testing, declare this file in the test directory.
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
	// DB is the representation of db connection settings.
	DB struct {
		// URI specifies the database instance address. A zero value
		// means there will a db connection error.
		URI string `env-required:"true" env:"DB_URI"`

		// Driver specifies the database's underlying driver.
		//
		// Default is postgres.
		Driver string `env:"DB_DRIVER" env-default:"postgres"`

		// DBname specifies the database's name to connect. A zero value
		// means there will a db connection error.
		DBname string `env-required:"true" env:"DB_NAME"`
	}

	// HTTP is the representation of http server settings
	HTTP struct {
		//Port specifies the listening port of the server.
		//
		// Default is 8080.
		Port string `env:"HTTP_PORT" env-default:"8080"`

		// ReadTimeout is the maximum duration for reading the entire
		// request, including the body. A zero or negative value means
		// there will be no timeout.
		//
		// Default is 5000us.
		ReadTimeout time.Duration `env:"HTTP_READ_TIMEOUT" env-default:"5000us"`

		// WriteTimeout is the maximum duration before timing out
		// writes of the response. A zero or negative value means
		// there will be no timeout.
		//
		// Default is 5000us.
		WriteTimeout time.Duration `env:"HTTP_WRITE_TIMEOUT" env-default:"5000us"`

		// ShutdownTimeout is the maximum duration before timing out
		// stops the running server. A zero or negative value means
		// there will be no timeout.
		//
		// Default is 3000us.
		ShutdownTimeout time.Duration `env:"HTTP_SHUTDOWN_TIMEOUT" env-default:"3000us"`
	}

	// Log is used for event logging configuration
	Log struct {
		// Level specifies the message importance level.
		//
		// Default is debug.
		Level string `env:"LOG_LEVEL" env-default:"debug"`
	}

	// Config holds all configuration structs, such as DB, HTTP, LOG
	Config struct {
		DB     DB
		HTTP   HTTP
		Logger Log
	}
)

var (
	profiles = map[string]string{
		"test":    "/test/.env",
		"default": ".env",
	}
)

func New(profile string) (Config, error) {
	var cfg Config

	cfgPath := fmt.Sprint(string(getRootDir()), profiles["default"])
	if v, ok := profiles[profile]; ok {
		cfgPath = fmt.Sprint(string(getRootDir()), v)
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
func (c *Config) Print() {
	if data, err := yaml.Marshal(c); err != nil {
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
