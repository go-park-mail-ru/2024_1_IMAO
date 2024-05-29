package metrics

import "time"

//go:generate mockgen -source=mock_gen.go -destination=mocks/mock.go

type DBMetrics interface {
	IncreaseErrors(funcName string)
	AddDuration(funcName string, duration time.Duration)
}
