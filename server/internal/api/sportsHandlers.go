package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/MedrekIT/sca-recrutation-task/server/internal/database"
)

func (cfg ApiConfig) addSportHandler(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1024)
	defer r.Body.Close()

	type addSport struct {
		Name string `json:"name"`
	}

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, []string{"INVALID_REQUEST", "Invalid sport data"}, fmt.Errorf("could not read request body: %w\n", err))
		return
	}

	var reqData addSport
	if err := json.Unmarshal(reqBody, &reqData); err != nil {
		errorResponse(w, http.StatusBadRequest, []string{"INVALID_REQUEST", "Invalid sport data"}, fmt.Errorf("could not decode request body: %w\n", err))
		return
	}

	var sport database.Sport
	if sport, err = cfg.Q.AddSport(r.Context(), reqData.Name); err != nil {
		errorResponse(w, http.StatusInternalServerError, []string{"DATABASE_ERROR", "Something went wrong"}, fmt.Errorf("could not add sport to the database - %w", err))
		return
	}

	successResponse(w, http.StatusCreated, struct {
		ID int32 `json:"id"`
	}{
		ID: sport.ID,
	})
}

func (cfg ApiConfig) getSportsHandler(w http.ResponseWriter, r *http.Request) {
	sports, err := cfg.Q.GetSports(r.Context())
	if err != nil {
		if err == sql.ErrNoRows {
			successResponse(w, http.StatusOK, []string{})
			return
		}
		errorResponse(w, http.StatusInternalServerError, []string{"DATABASE_ERROR", "Something went wrong"}, fmt.Errorf("could not get sports from the database - %w", err))
		return
	}

	successResponse(w, http.StatusOK, sports)
}
