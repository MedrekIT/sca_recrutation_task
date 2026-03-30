package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/MedrekIT/sca-recrutation-task/server/internal/database"
)

func (cfg ApiConfig) addCountryHandler(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1024)
	defer r.Body.Close()

	type addSport struct {
		CountryCode string `json:"country_code"`
		Name        string `json:"name"`
	}

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, []string{"INVALID_REQUEST", "Invalid country data"}, fmt.Errorf("could not read request body: %w\n", err))
		return
	}

	var reqData addSport
	if err := json.Unmarshal(reqBody, &reqData); err != nil {
		errorResponse(w, http.StatusBadRequest, []string{"INVALID_REQUEST", "Invalid country data"}, fmt.Errorf("could not decode request body: %w\n", err))
		return
	}

	if len(reqData.CountryCode) > 3 || reqData.CountryCode == "" || reqData.Name == "" {
		errorResponse(w, http.StatusBadRequest, []string{"INVALID_REQUEST", "Invalid country data"}, nil)
		return
	}

	newCountryParams := database.AddCountryParams{
		CountryCode: reqData.CountryCode,
		Name:        reqData.Name,
	}
	if _, err = cfg.Q.AddCountry(r.Context(), newCountryParams); err != nil {
		errorResponse(w, http.StatusInternalServerError, []string{"DATABASE_ERROR", "Something went wrong"}, fmt.Errorf("could not add country to the database - %w", err))
		return
	}

	successResponse(w, http.StatusCreated, struct {
		Message string `json:"message"`
	}{
		Message: "Country added successfully",
	})
}
