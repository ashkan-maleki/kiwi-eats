package config

import (
	"fmt"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"log"
)

var k = koanf.New(".")

type Config struct {
	AppName   string `koanf:"app_name"`
	Port      int    `koanf:"port"`
	JWTSecret string `koanf:"jwt_secret"`
	Database  struct {
		Host     string `koanf:"host"`
		Port     int    `koanf:"port"`
		User     string `koanf:"user"`
		Password string `koanf:"password"`
		DBName   string `koanf:"dbname"`
	} `koanf:"database"`
}

func LoadConfig() (*Config, error) {
	// Load YAML first (default values
	err := k.Load(file.Provider("config/config.dev.yaml"), yaml.Parser())
	if err != nil {
		log.Printf("Could not read YAML file: %v\n", err)
	}

	// Override with ENV variables
	k.Load(env.Provider("", ".", func(s string) string {
		return s
	}), nil)

	var cfg Config
	if err := k.Unmarshal("", &cfg); err != nil {
		return nil, fmt.Errorf("error unmarshalling config: %v", err)
	}
	return &cfg, nil
}
