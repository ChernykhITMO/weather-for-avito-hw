package usecase

import (
	"strings"

	"weather-service/internal/domain"
)

type WeatherService struct {
	repo domain.WeatherRepository
}

func NewWeatherService(repo domain.WeatherRepository) *WeatherService {
	return &WeatherService{repo: repo}
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

	return s.repo.GetByCity(city)
}
