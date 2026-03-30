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

	mu.HandleFunc("POST /api/sports", cfg.addSportHandler)
	mu.HandleFunc("POST /api/countries", cfg.addCountryHandler)
	mu.HandleFunc("POST /api/competitions", cfg.addCompetitionHandler)
	mu.HandleFunc("POST /api/competitions/{editionID}/stages", cfg.addStageHandler)
	mu.HandleFunc("POST /api/competitors", cfg.addCompetitorHandler)
	mu.HandleFunc("POST /api/venues", cfg.addVenueHandler)
	mu.HandleFunc("POST /api/events", cfg.createEventHandler)

	return mu
}
