package config

type RequestUUIDKey string
type LoggerKey string

const (
	RequestUUIDContextKey RequestUUIDKey = "requestUUID"
	LoggerContextKey      LoggerKey      = "logger"
)
