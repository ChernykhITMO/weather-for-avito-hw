package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"weather-service/internal/domain"
	"weather-service/internal/usecase"
)

type Handler struct {
	service *usecase.WeatherService
}

func NewHandler(service *usecase.WeatherService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", h.handleHealth)
	mux.HandleFunc("/weather", h.handleWeather)
	return mux
}

func (h *Handler) handleHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) handleWeather(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getWeather(w, r)
	case http.MethodPost:
		h.createWeather(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *Handler) getWeather(w http.ResponseWriter, r *http.Request) {
	city := r.URL.Query().Get("city")
	weather, err := h.service.GetByCity(city)
	if err != nil {
		h.writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, weather)
}

func (h *Handler) createWeather(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var weather domain.Weather
	if err := json.NewDecoder(r.Body).Decode(&weather); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}

	if err := h.service.Save(weather); err != nil {
		h.writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, weather)
}

func (h *Handler) writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrCityRequired):
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
	case errors.Is(err, domain.ErrNotFound):
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
	default:
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
