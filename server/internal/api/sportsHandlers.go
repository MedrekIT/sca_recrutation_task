package api

import (
	"database/sql"
	"fmt"
	"net/http"
)

func (cfg ApiConfig) getSportsHandler(w http.ResponseWriter, r *http.Request) {
	sports, err := cfg.Db.GetSports(r.Context())
	if err != nil {
		if err == sql.ErrNoRows {
			successResponse(w, http.StatusOK, struct {
				Sports []string `json:"sports"`
			}{
				Sports: []string{},
			})
			return
		}
		errorResponse(w, http.StatusInternalServerError, []string{"DATABASE_ERROR", "Something went wrong"}, fmt.Errorf("could not get sport from the database - %w", err))
		return
	}

	successResponse(w, http.StatusOK, struct {
		Sports []string `json:"sports"`
	}{
		Sports: sports,
	})
}
