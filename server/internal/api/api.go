package api

import (
	"database/sql"
	"net/http"

	"github.com/MedrekIT/sca-recrutation-task/server/internal/database"
)

type ApiConfig struct {
	Port string
	Db   *sql.DB
	Q    *database.Queries
}

func Routes(cfg ApiConfig) http.Handler {
	mu := http.NewServeMux()

	mu.HandleFunc("GET /api/events/{eventID}", cfg.getEventHandler)
	mu.HandleFunc("GET /api/events", cfg.getEventsHandler)
	mu.HandleFunc("GET /api/sports", cfg.getSportsHandler)
	mu.HandleFunc("GET /api/sports/{sportName}/competitions", cfg.getSportCompetitionsHandler)
	mu.HandleFunc("GET /api/sports/{sportName}/competitors", cfg.getSportCompetitorsHandler)
	mu.HandleFunc("GET /api/competitions/{editionID}/stages", cfg.getCompStagesHandler)
	mu.HandleFunc("GET /api/sports/{sportName}/venues", cfg.getSportVenuesHandler)

	// mu.HandleFunc("POST /api/sports", cfg.addSportsHandler)
	// mu.HandleFunc("POST /api/competitions", cfg.addCompetitionsHandler)
	// mu.HandleFunc("POST /api/competitors", cfg.addCompetitorsHandler)
	mu.HandleFunc("POST /api/events", cfg.createEventHandler)

	return mu
}
