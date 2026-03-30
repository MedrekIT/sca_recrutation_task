package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/MedrekIT/sca-recrutation-task/server/internal/database"
)

type Venue struct {
	CountryCode *string `json:"country_code,omitempty"`
	City        *string `json:"city,omitempty"`
	Place       *string `json:"place,omitempty"`
}

type Result struct {
	HomePoints    int16   `json:"home_points"`
	AwayPoints    int16   `json:"away_points"`
	Outcome       *string `json:"outcome,omitempty"`
	Winner        string  `json:"winner"`
	ResultDetails *string `json:"result_details,omitempty"`
}

type Event struct {
	EventID int32   `json:"event_id,omitempty"`
	Date    string  `json:"date"`
	Time    *string `json:"time,omitempty"`
	Venue   *Venue  `json:"venue,omitempty"`
	Status  string  `json:"status"`

	Competition string `json:"competition"`
	Season      string `json:"season"`
	Stage       string `json:"stage"`

	HomeCountry *string `json:"home_country,omitempty"`
	Home        *string `json:"home,omitempty"`
	AwayCountry *string `json:"away_country,omitempty"`
	Away        *string `json:"away,omitempty"`

	Result       Result  `json:"result"`
	EventDetails *string `json:"event_details,omitempty"`
}

func (cfg ApiConfig) getEventHandler(w http.ResponseWriter, r *http.Request) {
	eventID := r.PathValue("eventID")
	if eventID == "" {
		errorResponse(w, http.StatusBadRequest, []string{"INVALID_REQUEST", "Invalid event ID"}, nil)
	}

	intEventID, err := strconv.Atoi(eventID)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, []string{"INVALID_REQUEST", "Invalid event ID"}, fmt.Errorf("could not parse event ID from the URL path value: %w", err))
	}
	event, err := cfg.Db.GetEventByID(r.Context(), int32(intEventID))
	if err != nil {
		if err == sql.ErrNoRows {
			errorResponse(w, http.StatusNotFound, []string{"NOT_FOUND", "Event not found"}, fmt.Errorf("could not find event with given event ID: %w", err))
			return
		}
	}

	venueDate, venueTime := formatDateTime(event.VenueDate, event.VenueTime)

	winner := computeWinner(event.HomePoints, event.AwayPoints, event.Status, event.Outcome, event.ForfeitBy)

	var venue *Venue
	if event.CountryCode != nil || event.City != nil || event.PlaceName != nil {
		var countryPointer *string
		if event.CountryCode != nil {
			eventCountry := strings.TrimSpace(*event.CountryCode)
			countryPointer = &eventCountry
		}
		venue = &Venue{
			CountryCode: countryPointer,
			City:        event.City,
			Place:       event.PlaceName,
		}
	}

	eventRes := Event{
		Date:        venueDate,
		Time:        venueTime,
		Venue:       venue,
		Status:      event.Status,
		Competition: event.Competition,
		Season:      event.Season,
		Stage:       event.Stage,
		HomeCountry: event.HomeCountry,
		Home:        event.HomeCompetitor,
		AwayCountry: event.AwayCountry,
		Away:        event.AwayCompetitor,
		Result: Result{
			HomePoints:    event.HomePoints.Int16,
			AwayPoints:    event.AwayPoints.Int16,
			Outcome:       event.Outcome,
			Winner:        winner,
			ResultDetails: event.ResultDetails,
		},
		EventDetails: event.EventDetails,
	}

	successResponse(w, http.StatusOK, eventRes)
}

func (cfg ApiConfig) getEventsHandler(w http.ResponseWriter, r *http.Request) {
	filterQueries := r.URL.Query()

	var dateFilter sql.NullTime
	if filterQueries.Get("date") != "" {
		dateQuery, err := time.Parse("2006-01-02", filterQueries.Get("date"))
		if err != nil {
			errorResponse(w, http.StatusBadRequest, []string{"INVALID_REQUEST", "Invalid date"}, fmt.Errorf("could not parse date from the URL query: %w", err))
			return
		}
		dateFilter = sql.NullTime{
			Time:  dateQuery,
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
		errorResponse(w, http.StatusInternalServerError, []string{"DATABASE_ERROR", "Something went wrong"}, fmt.Errorf("couldn't get events from the database - %w", err))
		return
	}
	if events == nil {
		successResponse(w, http.StatusOK, []Event{})
		return
	}

	var eventsCalendar []Event
	for _, ev := range events {
		venueDate, venueTime := formatDateTime(ev.VenueDate, ev.VenueTime)

		winner := computeWinner(ev.HomePoints, ev.AwayPoints, ev.Status, ev.Outcome, ev.ForfeitBy)

		eventsCalendar = append(eventsCalendar, Event{
			EventID:     ev.EventID,
			Date:        venueDate,
			Time:        venueTime,
			Status:      ev.Status,
			Competition: ev.Competition,
			Season:      ev.Season,
			Stage:       ev.Stage,
			Home:        ev.HomeCompetitor,
			Away:        ev.AwayCompetitor,
			Result: Result{
				HomePoints: ev.HomePoints.Int16,
				AwayPoints: ev.AwayPoints.Int16,
				Winner:     winner,
			},
		})
	}

	successResponse(w, http.StatusOK, eventsCalendar)
}
