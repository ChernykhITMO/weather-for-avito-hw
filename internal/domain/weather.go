package domain

import "errors"

var (
	ErrCityRequired = errors.New("city is required")
	ErrNotFound     = errors.New("weather not found")
)

type Weather struct {
	City        string  `json:"city"`
	Temperature float64 `json:"temperature"`
	Condition   string  `json:"condition"`
}

type WeatherRepository interface {
	Save(weather Weather) error
	GetByCity(city string) (Weather, error)
}
