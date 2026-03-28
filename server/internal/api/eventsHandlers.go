package api

import (
	"database/sql"
	"fmt"
	"net/http"
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

	CompName string `json:"comp_name"`
	Season   string `json:"season"`
	Stage    string `json:"stage"`

	HomeComp string `json:"home_comp"`
	AwayComp string `json:"away_comp"`

	Result Result `json:"result"`
}

func (cfg ApiConfig) getEventsHandler(w http.ResponseWriter, r *http.Request) {
	events, err := cfg.Db.GetEvents(r.Context())
	if err != nil {
		if err == sql.ErrNoRows {
			successResponse(w, http.StatusOK, nil)
			return
		}
		errorResponse(w, http.StatusInternalServerError, []string{"DATABASE_ERROR", "Something went wrong"}, fmt.Errorf("couldn't get user from the database - %w", err))
		return
	}
	if events == nil {
		successResponse(w, http.StatusOK, nil)
		return
	}

	var eventsCalendar []Event
	for _, ev := range events {
		venueYear, venueMon, venueDay := ev.VenueDate.Date()
		venueDate := fmt.Sprintf("%d-%d-%d", venueYear, venueMon, venueDay)
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
		if ev.Outcome != nil && *ev.Outcome != "" {
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
			EventID:  ev.EventID,
			Date:     venueDate,
			Time:     &venueTime,
			Status:   ev.Status,
			CompName: ev.Competition,
			Season:   ev.Season,
			Stage:    ev.Stage,
			HomeComp: homeComp,
			AwayComp: awayComp,
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
