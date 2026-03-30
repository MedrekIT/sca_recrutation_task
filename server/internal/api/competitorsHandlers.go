package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/MedrekIT/sca-recrutation-task/server/internal/database"
)

func (cfg ApiConfig) addCompetitorHandler(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1024)
	defer r.Body.Close()

	type addCompetitor struct {
		Name        string  `json:"name"`
		CountryCode *string `json:"country_code"`
		SportID     int32   `json:"sport_id"`
	}

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, []string{"INVALID_REQUEST", "Invalid competitor data"}, fmt.Errorf("could not read request body: %w\n", err))
		return
	}

	var reqData addCompetitor
	if err := json.Unmarshal(reqBody, &reqData); err != nil {
		errorResponse(w, http.StatusBadRequest, []string{"INVALID_REQUEST", "Invalid competitor data"}, fmt.Errorf("could not decode request body: %w\n", err))
		return
	}

	if reqData.CountryCode != nil && (len(*reqData.CountryCode) > 3 || *reqData.CountryCode == "") {
		errorResponse(w, http.StatusBadRequest, []string{"INVALID_REQUEST", "Invalid country code"}, nil)
		return
	}

	if reqData.CountryCode != nil {
		if _, err := cfg.Q.GetCountryByCode(r.Context(), *reqData.CountryCode); err != nil {
			if err == sql.ErrNoRows {
				errorResponse(w, http.StatusBadRequest, []string{"INVALID_REQUEST", "Invalid country code"}, fmt.Errorf("could not get country from the database - %w", err))
				return
			}
			errorResponse(w, http.StatusInternalServerError, []string{"DATABASE_ERROR", "Something went wrong"}, fmt.Errorf("could not get country from the database - %w", err))
			return
		}
	}

	if _, err := cfg.Q.GetSportByID(r.Context(), reqData.SportID); err != nil {
		if err == sql.ErrNoRows {
			errorResponse(w, http.StatusBadRequest, []string{"INVALID_REQUEST", "Invalid sport ID"}, fmt.Errorf("could not get sport from the database - %w", err))
			return
		}
		errorResponse(w, http.StatusInternalServerError, []string{"DATABASE_ERROR", "Something went wrong"}, fmt.Errorf("could not get sport from the database - %w", err))
		return
	}

	newCompetitorParams := database.AddCompetitorParams{
		Name:    reqData.Name,
		Country: reqData.CountryCode,
		SportID: reqData.SportID,
	}
	competitorID, err := cfg.Q.AddCompetitor(r.Context(), newCompetitorParams)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, []string{"DATABASE_ERROR", "Something went wrong"}, fmt.Errorf("could not add competition to the database: %w", err))
		return
	}

	successResponse(w, http.StatusCreated, struct {
		CompetitorID int32 `json:"competitor_id"`
	}{
		CompetitorID: competitorID,
	})
}

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
