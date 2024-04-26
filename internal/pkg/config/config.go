package config

import (
	"time"
)

type RequestUUIDKey string
type LoggerKey string
type SessionKey string

const (
	RequestUUIDContextKey RequestUUIDKey = "requestUUID"
	LoggerContextKey      LoggerKey      = "logger"
	SessionContextKey     SessionKey     = "session"
)

type CsrfConfig struct {
	CsrfCookie   string        `yaml:"csrf_cookie"`
	CSRFLifeTime time.Duration `yaml:"csrf_lifetime"`
}
