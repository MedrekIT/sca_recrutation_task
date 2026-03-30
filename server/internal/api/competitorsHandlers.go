package api

import (
	"database/sql"
	"fmt"
	"net/http"
)

func (cfg ApiConfig) getSportCompetitorsHandler(w http.ResponseWriter, r *http.Request) {
	sportName := r.PathValue("sportName")
	if sportName == "" {
		errorResponse(w, http.StatusBadRequest, []string{"INVALID_REQUEST", "Invalid sport name"}, nil)
		return
	}

	competitors, err := cfg.Q.GetCompetitorsForSport(r.Context(), sportName)
	if err != nil {
		if err == sql.ErrNoRows {
			successResponse(w, http.StatusOK, []string{})
			return
		}
		errorResponse(w, http.StatusInternalServerError, []string{"DATABASE_ERROR", "Something went wrong"}, fmt.Errorf("could not get competitors from the database - %w", err))
		return
	}

	successResponse(w, http.StatusOK, competitors)
}
