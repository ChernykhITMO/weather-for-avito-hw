package usecase

import (
	"errors"
	"testing"
	"time"

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

func TestWeatherService_GetByCity_SavesHistoryOnSuccess(t *testing.T) {
	repo := memoryrepo.NewWeatherRepository()
	service := NewWeatherService(repo)
	fixedTime := time.Date(2026, 4, 9, 15, 4, 5, 0, time.UTC)
	service.now = func() time.Time { return fixedTime }

	input := domain.Weather{
		City:        "Moscow",
		Temperature: 18.5,
		Condition:   "Cloudy",
	}

	if err := service.Save(input); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	got, err := service.GetByCity("Moscow")
	if err != nil {
		t.Fatalf("GetByCity() error = %v", err)
	}

	if got.City != input.City {
		t.Fatalf("GetByCity().City = %q, want %q", got.City, input.City)
	}

	history, err := service.GetHistory()
	if err != nil {
		t.Fatalf("GetHistory() error = %v", err)
	}

	if len(history) != 1 {
		t.Fatalf("len(GetHistory()) = %d, want 1", len(history))
	}

	if history[0].RequestedAt != fixedTime {
		t.Fatalf("GetHistory()[0].RequestedAt = %v, want %v", history[0].RequestedAt, fixedTime)
	}

	if history[0].City != input.City || history[0].Condition != input.Condition || history[0].Temperature != input.Temperature {
		t.Fatalf("GetHistory()[0] = %+v, want weather fields copied from %+v", history[0], input)
	}
}

func TestWeatherService_GetByCity_DoesNotSaveHistoryOnFailure(t *testing.T) {
	service := NewWeatherService(memoryrepo.NewWeatherRepository())

	_, err := service.GetByCity("Unknown")
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("GetByCity() error = %v, want %v", err, domain.ErrNotFound)
	}

	history, err := service.GetHistory()
	if err != nil {
		t.Fatalf("GetHistory() error = %v", err)
	}

	if len(history) != 0 {
		t.Fatalf("len(GetHistory()) = %d, want 0", len(history))
	}
}

type saveHistoryFailingRepo struct {
	weather domain.Weather
	err     error
}

func (r saveHistoryFailingRepo) Save(weather domain.Weather) error {
	return nil
}

func (r saveHistoryFailingRepo) GetByCity(city string) (domain.Weather, error) {
	return r.weather, nil
}

func (r saveHistoryFailingRepo) SaveHistory(record domain.HistoryRecord) error {
	return r.err
}

func (r saveHistoryFailingRepo) GetHistory() ([]domain.HistoryRecord, error) {
	return nil, nil
}

func TestWeatherService_GetByCity_IgnoresSaveHistoryError(t *testing.T) {
	expected := domain.Weather{
		City:        "Moscow",
		Temperature: 18.5,
		Condition:   "Cloudy",
	}
	service := NewWeatherService(saveHistoryFailingRepo{
		weather: expected,
		err:     errors.New("history unavailable"),
	})

	got, err := service.GetByCity("Moscow")
	if err != nil {
		t.Fatalf("GetByCity() error = %v, want nil", err)
	}

	if got != expected {
		t.Fatalf("GetByCity() = %+v, want %+v", got, expected)
	}
}
