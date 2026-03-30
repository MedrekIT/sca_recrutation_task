package api

import (
	"net/http"

	"github.com/MedrekIT/sca-recrutation-task/server/internal/database"
)

type ApiConfig struct {
	Port string
	Db   *database.Queries
}

func Routes(cfg ApiConfig) http.Handler {
	mu := http.NewServeMux()

	mu.HandleFunc("GET /api/events/{eventID}", cfg.getEventHandler)
	mu.HandleFunc("GET /api/events", cfg.getEventsHandler)
	mu.HandleFunc("GET /api/sports", cfg.getSportsHandler)

	return mu
}
