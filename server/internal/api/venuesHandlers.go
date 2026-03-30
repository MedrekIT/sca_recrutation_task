package api

import (
	"database/sql"
	"fmt"
	"net/http"
)

type Venues struct {
	VenueID int32  `json:"venue_id"`
	Name    string `json:"name"`
}

func (cfg ApiConfig) getSportVenuesHandler(w http.ResponseWriter, r *http.Request) {
	sportName := r.PathValue("sportName")
	if sportName == "" {
		errorResponse(w, http.StatusBadRequest, []string{"INVALID_REQUEST", "Invalid sport name"}, nil)
		return
	}

	venues, err := cfg.Q.GetVenuesForSport(r.Context(), sportName)
	if err != nil {
		if err == sql.ErrNoRows {
			successResponse(w, http.StatusOK, []string{})
			return
		}
		errorResponse(w, http.StatusInternalServerError, []string{"DATABASE_ERROR", "Something went wrong"}, fmt.Errorf("could not get venues from the database - %w", err))
		return
	}

	var venuesRes []Venues
	for _, venue := range venues {
		venuesRes = append(venuesRes, Venues{
			VenueID: venue.ID,
			Name:    venue.Name,
		})
	}

	successResponse(w, http.StatusOK, venuesRes)
}
