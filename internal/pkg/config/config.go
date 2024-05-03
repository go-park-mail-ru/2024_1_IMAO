package config

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"time"
)

type RequestUUIDKey string
type LoggerKey string
type SessionKey string

const (
	RequestUUIDContextKey RequestUUIDKey = "requestUUID"
	LoggerContextKey      LoggerKey      = "logger"
	SessionContextKey     SessionKey     = "session"
	cfgPath                              = "./internal/pkg/config/config.yaml"
)

type CsrfConfig struct {
	CsrfCookie   string        `yaml:"csrf_cookie"`
	CSRFLifeTime time.Duration `yaml:"csrf_lifetime"`
}

type ServerConfig struct {
	Host               string `yaml:"host"`
	Port               string `yaml:"port"`
	AuthServicePort    string `yaml:"auth_service_port"`
	ProfileServicePort string `yaml:"profile_service_port"`
}

type Config struct {
	Server ServerConfig `yaml:"server"`
}

func ReadConfig() *Config {
	cfg := &Config{}

	file, err := os.Open(cfgPath)
	if err != nil {
		log.Println("Something went wrong while opening config file:", err)
		return nil
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(cfg); err != nil {
		log.Println("Something went wrong while reading config file:", err)
		return nil
	}

	log.Println("Successfully opened config")
	return cfg
}
