package usecase

import (
	"errors"
	"testing"

	"weather-service/internal/domain"
	memoryrepo "weather-service/internal/repository/memory"
)

func TestWeatherService_SaveAndGetByCity(t *testing.T) {
	service := NewWeatherService(memoryrepo.NewWeatherRepository())
	input := domain.Weather{
		City:        "Moscow",
		Temperature: 17.5,
		Condition:   "Cloudy",
	}

	if err := service.Save(input); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	got, err := service.GetByCity("moscow")
	if err != nil {
		t.Fatalf("GetByCity() error = %v", err)
	}

	if got.City != input.City {
		t.Fatalf("GetByCity().City = %q, want %q", got.City, input.City)
	}
}

func TestWeatherService_Save_EmptyCity(t *testing.T) {
	service := NewWeatherService(memoryrepo.NewWeatherRepository())

	err := service.Save(domain.Weather{})
	if !errors.Is(err, domain.ErrCityRequired) {
		t.Fatalf("Save() error = %v, want %v", err, domain.ErrCityRequired)
	}
}
