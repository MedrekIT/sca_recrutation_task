package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/MedrekIT/sca-recrutation-task/server/internal/database"
)

type Stage struct {
	StageID int32  `json:"stage_id"`
	Name    string `json:"name"`
}

type Competition struct {
	EditionID int32  `json:"edition_id"`
	Label     string `json:"label"`
}

func (cfg ApiConfig) addStageHandler(w http.ResponseWriter, r *http.Request) {
	editionID := r.PathValue("editionID")
	if editionID == "" {
		errorResponse(w, http.StatusBadRequest, []string{"INVALID_REQUEST", "Invalid edition ID"}, nil)
		return
	}
	intEditionID, err := strconv.Atoi(editionID)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, []string{"INVALID_REQUEST", "Invalid edition ID"}, fmt.Errorf("could not parse edition ID from the URL path value: %w", err))
		return
	}
	edition, err := cfg.Q.GetEditionByID(r.Context(), int32(intEditionID))
	if err != nil {
		if err == sql.ErrNoRows {
			errorResponse(w, http.StatusBadRequest, []string{"INVALID_REQUEST", "Invalid edition ID"}, fmt.Errorf("could not find edition with given edition ID: %w", err))
			return
		}
		errorResponse(w, http.StatusInternalServerError, []string{"DATABASE_ERROR", "Something went wrong"}, fmt.Errorf("could not get edition from the database: %w", err))
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 1024)
	defer r.Body.Close()

	type addStage struct {
		Name string `json:"name"`
	}

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, []string{"INVALID_REQUEST", "Invalid stage data"}, fmt.Errorf("could not read request body: %w\n", err))
		return
	}

	var reqData addStage
	if err := json.Unmarshal(reqBody, &reqData); err != nil {
		errorResponse(w, http.StatusBadRequest, []string{"INVALID_REQUEST", "Invalid stage data"}, fmt.Errorf("could not decode request body: %w\n", err))
		return
	}

	newStageParams := database.AddStageParams{
		Name:      reqData.Name,
		EditionID: edition.ID,
	}
	stageID, err := cfg.Q.AddStage(r.Context(), newStageParams)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, []string{"DATABASE_ERROR", "Something went wrong"}, fmt.Errorf("could not add stage to the database: %w", err))
		return
	}

	successResponse(w, http.StatusCreated, struct {
		ID int32 `json:"id"`
	}{
		ID: stageID,
	})
}

func (cfg ApiConfig) addCompetitionHandler(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1024)
	defer r.Body.Close()

	type addCompetition struct {
		CompetitionName string `json:"competition_name"`
		SportID         int32  `json:"sport_id"`
		Season          string `json:"season"`
	}

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, []string{"INVALID_REQUEST", "Invalid event data"}, fmt.Errorf("could not read request body: %w\n", err))
		return
	}

	var reqData addCompetition
	if err := json.Unmarshal(reqBody, &reqData); err != nil {
		errorResponse(w, http.StatusBadRequest, []string{"INVALID_REQUEST", "Invalid event data"}, fmt.Errorf("could not decode request body: %w\n", err))
		return
	}

	if _, err := cfg.Q.GetSportByID(r.Context(), reqData.SportID); err != nil {
		if err == sql.ErrNoRows {
			errorResponse(w, http.StatusBadRequest, []string{"INVALID_REQUEST", "Invalid sport ID"}, fmt.Errorf("could not get sport from the database - %w", err))
			return
		}
		errorResponse(w, http.StatusInternalServerError, []string{"DATABASE_ERROR", "Something went wrong"}, fmt.Errorf("could not get sport from the database - %w", err))
		return
	}

	tx, err := cfg.Db.BeginTx(r.Context(), nil)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, []string{"DATABASE_ERROR", "Something went wrong"}, fmt.Errorf("could not start database transaction: %w", err))
		return
	}
	defer tx.Rollback()

	qtx := cfg.Q.WithTx(tx)

	newCompetitionParams := database.AddCompetitionParams{
		Name:    reqData.CompetitionName,
		SportID: reqData.SportID,
	}
	competitionID, err := qtx.AddCompetition(r.Context(), newCompetitionParams)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, []string{"DATABASE_ERROR", "Something went wrong"}, fmt.Errorf("could not add competition to the database: %w", err))
		return
	}

	newEditionParams := database.AddEditionParams{
		Season:        reqData.Season,
		CompetitionID: competitionID,
	}
	editionID, err := qtx.AddEdition(r.Context(), newEditionParams)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, []string{"DATABASE_ERROR", "Something went wrong"}, fmt.Errorf("could not add edition to the database: %w", err))
		return
	}

	if err := tx.Commit(); err != nil {
		errorResponse(w, http.StatusInternalServerError, []string{"DATABASE_ERROR", "Something went wrong"}, fmt.Errorf("could not commit database transaction: %w", err))
		return
	}

	successResponse(w, http.StatusCreated, struct {
		ID int32 `json:"id"`
	}{
		ID: editionID,
	})
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
