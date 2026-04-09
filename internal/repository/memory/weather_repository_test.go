package memory

import (
	"fmt"
	"testing"
	"time"

	"weather-service/internal/domain"
)

func TestWeatherRepository_SaveHistory_NewestFirstAndLimitedToTen(t *testing.T) {
	repo := NewWeatherRepository()
	base := time.Date(2026, 4, 9, 15, 0, 0, 0, time.UTC)

	for i := 0; i < 12; i++ {
		record := domain.HistoryRecord{
			City:        fmt.Sprintf("City-%02d", i),
			Temperature: float64(i),
			Condition:   "Sunny",
			RequestedAt: base.Add(time.Duration(i) * time.Minute),
		}

		if err := repo.SaveHistory(record); err != nil {
			t.Fatalf("SaveHistory() error = %v", err)
		}
	}

	got, err := repo.GetHistory()
	if err != nil {
		t.Fatalf("GetHistory() error = %v", err)
	}

	if len(got) != 10 {
		t.Fatalf("len(GetHistory()) = %d, want 10", len(got))
	}

	if got[0].City != "City-11" {
		t.Fatalf("GetHistory()[0].City = %q, want %q", got[0].City, "City-11")
	}

	if got[len(got)-1].City != "City-02" {
		t.Fatalf("GetHistory()[%d].City = %q, want %q", len(got)-1, got[len(got)-1].City, "City-02")
	}
}
