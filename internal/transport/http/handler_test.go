package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"weather-service/internal/domain"
	memoryrepo "weather-service/internal/repository/memory"
	"weather-service/internal/usecase"
)

func TestHandler_CreateAndGetWeather(t *testing.T) {
	handler := NewHandler(usecase.NewWeatherService(memoryrepo.NewWeatherRepository()))
	server := handler.Routes()

	postReq := httptest.NewRequest(http.MethodPost, "/weather", bytes.NewBufferString(`{"city":"Kazan","temperature":21.3,"condition":"Sunny"}`))
	postRec := httptest.NewRecorder()
	server.ServeHTTP(postRec, postReq)

	if postRec.Code != http.StatusCreated {
		t.Fatalf("POST /weather status = %d, want %d", postRec.Code, http.StatusCreated)
	}

	getReq := httptest.NewRequest(http.MethodGet, "/weather?city=Kazan", nil)
	getRec := httptest.NewRecorder()
	server.ServeHTTP(getRec, getReq)

	if getRec.Code != http.StatusOK {
		t.Fatalf("GET /weather status = %d, want %d", getRec.Code, http.StatusOK)
	}
}

func TestHandler_GetHistory_EmptyList(t *testing.T) {
	handler := NewHandler(usecase.NewWeatherService(memoryrepo.NewWeatherRepository()))
	server := handler.Routes()

	req := httptest.NewRequest(http.MethodGet, "/history", nil)
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("GET /history status = %d, want %d", rec.Code, http.StatusOK)
	}

	var got []domain.HistoryRecord
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	if len(got) != 0 {
		t.Fatalf("len(history) = %d, want 0", len(got))
	}
}

func TestHandler_GetHistory_ReturnsNewestFirst(t *testing.T) {
	handler := NewHandler(usecase.NewWeatherService(memoryrepo.NewWeatherRepository()))
	server := handler.Routes()

	for _, payload := range []string{
		`{"city":"Moscow","temperature":18.5,"condition":"Cloudy"}`,
		`{"city":"Kazan","temperature":21.3,"condition":"Sunny"}`,
		`{"city":"Sochi","temperature":24.1,"condition":"Clear"}`,
	} {
		postReq := httptest.NewRequest(http.MethodPost, "/weather", bytes.NewBufferString(payload))
		postRec := httptest.NewRecorder()
		server.ServeHTTP(postRec, postReq)
		if postRec.Code != http.StatusCreated {
			t.Fatalf("POST /weather status = %d, want %d", postRec.Code, http.StatusCreated)
		}
	}

	for _, city := range []string{"Moscow", "Kazan", "Sochi"} {
		getReq := httptest.NewRequest(http.MethodGet, "/weather?city="+city, nil)
		getRec := httptest.NewRecorder()
		server.ServeHTTP(getRec, getReq)
		if getRec.Code != http.StatusOK {
			t.Fatalf("GET /weather?city=%s status = %d, want %d", city, getRec.Code, http.StatusOK)
		}
	}

	historyReq := httptest.NewRequest(http.MethodGet, "/history", nil)
	historyRec := httptest.NewRecorder()
	server.ServeHTTP(historyRec, historyReq)

	if historyRec.Code != http.StatusOK {
		t.Fatalf("GET /history status = %d, want %d", historyRec.Code, http.StatusOK)
	}

	var got []domain.HistoryRecord
	if err := json.Unmarshal(historyRec.Body.Bytes(), &got); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	if len(got) != 3 {
		t.Fatalf("len(history) = %d, want 3", len(got))
	}

	wantOrder := []string{"Sochi", "Kazan", "Moscow"}
	for i, wantCity := range wantOrder {
		if got[i].City != wantCity {
			t.Fatalf("history[%d].City = %q, want %q", i, got[i].City, wantCity)
		}
		if got[i].RequestedAt.IsZero() {
			t.Fatalf("history[%d].RequestedAt is zero", i)
		}
	}
}
