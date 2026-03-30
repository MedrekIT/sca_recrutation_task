package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
)

type Stage struct {
	StageID int32  `json:"stage_id"`
	Name    string `json:"name"`
}

type Competition struct {
	EditionID int32  `json:"edition_id"`
	Label     string `json:"label"`
}

func (cfg ApiConfig) getCompStagesHandler(w http.ResponseWriter, r *http.Request) {
	editionID := r.PathValue("editionID")
	if editionID == "" {
		errorResponse(w, http.StatusBadRequest, []string{"INVALID_REQUEST", "Invalid competition edition ID"}, nil)
		return
	}
	intEditionID, err := strconv.Atoi(editionID)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, []string{"INVALID_REQUEST", "Invalid competition edition ID"}, fmt.Errorf("could not parse edition ID from the URL path value: %w", err))
		return
	}

	stages, err := cfg.Q.GetStagesForCompEdition(r.Context(), int32(intEditionID))
	if err != nil {
		if err == sql.ErrNoRows {
			successResponse(w, http.StatusOK, []string{})
			return
		}
		errorResponse(w, http.StatusInternalServerError, []string{"DATABASE_ERROR", "Something went wrong"}, fmt.Errorf("could not get stages from the database - %w", err))
		return
	}

	var stagesRes []Stage
	for _, stage := range stages {
		stagesRes = append(stagesRes, Stage{
			StageID: stage.ID,
			Name:    stage.Name,
		})
	}

	successResponse(w, http.StatusOK, stagesRes)
}

func (cfg ApiConfig) getSportCompetitionsHandler(w http.ResponseWriter, r *http.Request) {
	sportName := r.PathValue("sportName")
	if sportName == "" {
		errorResponse(w, http.StatusBadRequest, []string{"INVALID_REQUEST", "Invalid sport name"}, nil)
		return
	}

	competitions, err := cfg.Q.GetCompEditionsForSport(r.Context(), sportName)
	if err != nil {
		if err == sql.ErrNoRows {
			successResponse(w, http.StatusOK, []string{})
			return
		}
		errorResponse(w, http.StatusInternalServerError, []string{"DATABASE_ERROR", "Something went wrong"}, fmt.Errorf("could not get competitions from the database - %w", err))
		return
	}

	var competitionsRes []Competition
	for _, comp := range competitions {
		competitionsRes = append(competitionsRes, Competition{
			EditionID: comp.ID,
			Label:     fmt.Sprintf("%s %s", comp.Name, comp.Season),
		})
	}

	successResponse(w, http.StatusOK, competitionsRes)
}
