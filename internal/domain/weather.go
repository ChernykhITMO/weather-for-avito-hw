package domain

import (
	"errors"
	"time"
)

var (
	ErrCityRequired = errors.New("city is required")
	ErrNotFound     = errors.New("weather not found")
)

type Weather struct {
	City        string  `json:"city"`
	Temperature float64 `json:"temperature"`
	Condition   string  `json:"condition"`
}

type HistoryRecord struct {
	City        string    `json:"city"`
	Temperature float64   `json:"temperature"`
	Condition   string    `json:"condition"`
	RequestedAt time.Time `json:"requested_at"`
}

type WeatherRepository interface {
	Save(weather Weather) error
	GetByCity(city string) (Weather, error)
	SaveHistory(record HistoryRecord) error
	GetHistory() ([]HistoryRecord, error)
}
