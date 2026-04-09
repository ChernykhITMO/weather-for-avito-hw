package http

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

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
