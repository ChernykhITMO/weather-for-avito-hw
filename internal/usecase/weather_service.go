package usecase

import (
	"strings"
	"time"

	"weather-service/internal/domain"
)

type WeatherService struct {
	repo domain.WeatherRepository
	now  func() time.Time
}

func NewWeatherService(repo domain.WeatherRepository) *WeatherService {
	return &WeatherService{
		repo: repo,
		now:  time.Now,
	}
}

func (s *WeatherService) Save(weather domain.Weather) error {
	weather.City = strings.TrimSpace(weather.City)
	weather.Condition = strings.TrimSpace(weather.Condition)
	if weather.City == "" {
		return domain.ErrCityRequired
	}

	return s.repo.Save(weather)
}

func (s *WeatherService) GetByCity(city string) (domain.Weather, error) {
	city = strings.TrimSpace(city)
	if city == "" {
		return domain.Weather{}, domain.ErrCityRequired
	}

	weather, err := s.repo.GetByCity(city)
	if err != nil {
		return domain.Weather{}, err
	}

	record := domain.HistoryRecord{
		City:        weather.City,
		Temperature: weather.Temperature,
		Condition:   weather.Condition,
		RequestedAt: s.now().UTC(),
	}
	if err := s.repo.SaveHistory(record); err != nil {
		return domain.Weather{}, err
	}

	return weather, nil
}

func (s *WeatherService) GetHistory() ([]domain.HistoryRecord, error) {
	return s.repo.GetHistory()
}
