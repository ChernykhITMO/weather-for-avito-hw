package memory

import (
	"strings"
	"sync"

	"weather-service/internal/domain"
)

type WeatherRepository struct {
	mu      sync.RWMutex
	weather map[string]domain.Weather
}

func NewWeatherRepository() *WeatherRepository {
	return &WeatherRepository{
		weather: make(map[string]domain.Weather),
	}
}

func (r *WeatherRepository) Save(weather domain.Weather) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.weather[normalizeCity(weather.City)] = weather
	return nil
}

func (r *WeatherRepository) GetByCity(city string) (domain.Weather, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	weather, ok := r.weather[normalizeCity(city)]
	if !ok {
		return domain.Weather{}, domain.ErrNotFound
	}

	return weather, nil
}

func normalizeCity(city string) string {
	return strings.ToLower(strings.TrimSpace(city))
}
