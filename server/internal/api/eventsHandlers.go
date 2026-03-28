package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/MedrekIT/sca-recrutation-task/server/internal/database"
)

type Result struct {
	HomePoints int16   `json:"home_points"`
	AwayPoints int16   `json:"away_points"`
	Outcome    *string `json:"outcome,omitempty"`
	Winner     string  `json:"winner"`
}

type Event struct {
	EventID int32   `json:"event_id"`
	Date    string  `json:"date"`
	Time    *string `json:"time,omitempty"`
	Status  string  `json:"status"`

	Competition string `json:"competition"`
	Season      string `json:"season"`
	Stage       string `json:"stage"`

	Home string `json:"home"`
	Away string `json:"away"`

	Result Result `json:"result"`
}

func (cfg ApiConfig) getEventHandler(w http.ResponseWriter, r *http.Request) {
}

func (cfg ApiConfig) getEventsHandler(w http.ResponseWriter, r *http.Request) {
	filterQueries := r.URL.Query()

	var dateFilter sql.NullTime
	if filterQueries.Get("date") != "" {
		dateQuery, err := time.Parse("2006-01-02", filterQueries.Get("date"))
		if err != nil {
			errorResponse(w, http.StatusBadRequest, []string{"INVALID_REQUEST", "Invalid date"}, fmt.Errorf("couldn't parse date from the URL query: %w", err))
			return
		}
		dateFilter = sql.NullTime{
			Time: dateQuery,
			Valid: true,
		}
	}
	var sportFilter *string
	if filterQueries.Get("sport") != "" {
		sportQuery := filterQueries.Get("sport")
		sportFilter = &sportQuery
	}
	getEventsParams := database.GetEventsParams{
		DateFilter:  dateFilter,
		SportFilter: sportFilter,
	}
	events, err := cfg.Db.GetEvents(r.Context(), getEventsParams)
	if err != nil {
		if err == sql.ErrNoRows {
			successResponse(w, http.StatusOK, []Event{})
			return
		}
		errorResponse(w, http.StatusInternalServerError, []string{"DATABASE_ERROR", "Something went wrong"}, fmt.Errorf("couldn't get user from the database - %w", err))
		return
	}
	if events == nil {
		successResponse(w, http.StatusOK, []Event{})
		return
	}

	var eventsCalendar []Event
	for _, ev := range events {
		venueYear, venueMon, venueDay := ev.VenueDate.Date()
		venueDate := fmt.Sprintf("%d/%02d/%02d", venueYear, venueMon, venueDay)
		venueHour, venueMin, _ := ev.VenueTime.Time.Clock()
		venueTime := fmt.Sprintf("%02d:%02d", venueHour, venueMin)

		var winner string
		if ev.HomePoints.Valid && ev.AwayPoints.Valid {
			if ev.HomePoints.Int16 > ev.AwayPoints.Int16 {
				winner = "home"
			} else if ev.HomePoints.Int16 < ev.AwayPoints.Int16 {
				winner = "away"
			}
		}
		status := ev.Status
		if ev.Outcome != nil && *ev.Outcome != "" {
			status = "finished"
			if *ev.Outcome == "forfeit" {
				if ev.ForfeitBy != nil && *ev.ForfeitBy == "home" {
					winner = "away"
				}
				if ev.ForfeitBy != nil && *ev.ForfeitBy == "away" {
					winner = "home"
				}
			} else {
				winner = ""
			}
		}

		var homeComp, awayComp string
		if ev.HomeCompetitor != nil {
			homeComp = *ev.HomeCompetitor
		}
		if ev.AwayCompetitor != nil {
			awayComp = *ev.AwayCompetitor
		}

		eventsCalendar = append(eventsCalendar, Event{
			EventID:     ev.EventID,
			Date:        venueDate,
			Time:        &venueTime,
			Status:      status,
			Competition: ev.Competition,
			Season:      ev.Season,
			Stage:       ev.Stage,
			Home:        homeComp,
			Away:        awayComp,
			Result: Result{
				HomePoints: ev.HomePoints.Int16,
				AwayPoints: ev.AwayPoints.Int16,
				Outcome:    ev.Outcome,
				Winner:     winner,
			},
		})
	}

	successResponse(w, http.StatusOK, eventsCalendar)
}
