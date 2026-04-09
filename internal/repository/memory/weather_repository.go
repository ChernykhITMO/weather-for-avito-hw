package memory

import (
	"strings"
	"sync"

	"weather-service/internal/domain"
)

type WeatherRepository struct {
	mu      sync.RWMutex
	weather map[string]domain.Weather
	history []domain.HistoryRecord
}

func NewWeatherRepository() *WeatherRepository {
	return &WeatherRepository{
		weather: make(map[string]domain.Weather),
		history: make([]domain.HistoryRecord, 0, 10),
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

func (r *WeatherRepository) SaveHistory(record domain.HistoryRecord) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.history = append([]domain.HistoryRecord{record}, r.history...)
	if len(r.history) > 10 {
		r.history = r.history[:10]
	}

	return nil
}

func (r *WeatherRepository) GetHistory() ([]domain.HistoryRecord, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	history := make([]domain.HistoryRecord, len(r.history))
	copy(history, r.history)
	return history, nil
}

func normalizeCity(city string) string {
	return strings.ToLower(strings.TrimSpace(city))
}
