package main

import (
	"log"
	"net/http"

	memoryrepo "weather-service/internal/repository/memory"
	httptransport "weather-service/internal/transport/http"
	"weather-service/internal/usecase"
)

func main() {
	repo := memoryrepo.NewWeatherRepository()
	service := usecase.NewWeatherService(repo)
	handler := httptransport.NewHandler(service)

	server := &http.Server{
		Addr:    ":8080",
		Handler: handler.Routes(),
	}

	log.Printf("weather-service listening on %s", server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server failed: %v", err)
	}
}
