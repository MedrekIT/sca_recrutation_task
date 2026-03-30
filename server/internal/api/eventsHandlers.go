package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"
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
	ForteitBy     *string `json:"forfeit_by,omitempty"`
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

func (cfg ApiConfig) createEventHandler(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 2048)
	defer r.Body.Close()

	type createEvent struct {
		Status    string     `json:"status"`
		VenueDate time.Time  `json:"venue_date"`
		VenueTime *time.Time `json:"venue_time,omitempty"`
		VenueID   *int32     `json:"venue_id,omitempty"`
		HomeName  *string    `json:"home_name,omitempty"`
		AwayName  *string    `json:"away_name,omitempty"`
		StageID   int32      `json:"stage_id"`
		Result    *Result    `json:"result,omitempty"`
		Details   *string    `json:"details,omitempty"`
	}

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, []string{"INVALID_REQUEST", "Invalid event data"}, fmt.Errorf("could not read request body: %w\n", err))
		return
	}

	var reqData createEvent
	if err := json.Unmarshal(reqBody, &reqData); err != nil {
		errorResponse(w, http.StatusBadRequest, []string{"INVALID_REQUEST", "Invalid event data"}, fmt.Errorf("could not decode request body: %w\n", err))
		return
	}

	if !slices.Contains(availableStatus, reqData.Status) {
		errorResponse(w, http.StatusBadRequest, []string{"INVALID_REQUEST", "Invalid event status"}, nil)
		return
	}
	if reqData.VenueDate.IsZero() {
		errorResponse(w, http.StatusBadRequest, []string{"INVALID_REQUEST", "Invalid event date"}, nil)
		return
	}
	if (reqData.HomeName != nil && reqData.AwayName != nil) && (*reqData.HomeName == *reqData.AwayName) {
		errorResponse(w, http.StatusBadRequest, []string{"INVALID_REQUEST", "Invalid teams names"}, nil)
		return
	}

	var sportID int32
	var venueID sql.NullInt32
	if reqData.VenueID != nil {
		venue, err := cfg.Q.GetVenueByID(r.Context(), *reqData.VenueID)
		if err != nil {
			if err == sql.ErrNoRows {
				errorResponse(w, http.StatusBadRequest, []string{"INVALID_REQUEST", "Invalid venue"}, fmt.Errorf("could not find venue with given venue ID: %w", err))
				return
			}
			errorResponse(w, http.StatusInternalServerError, []string{"DATABASE_ERROR", "Something went wrong"}, fmt.Errorf("could not get venue from the database: %w", err))
			return
		}

		sportID = venue.SportID
		venueID = sql.NullInt32{
			Int32: venue.ID,
			Valid: true,
		}
	}

	var homeID, awayID sql.NullInt32
	if reqData.HomeName != nil {
		home, err := cfg.Q.GetCompetitorByName(r.Context(), *reqData.HomeName)
		if err != nil {
			if err == sql.ErrNoRows {
				errorResponse(w, http.StatusBadRequest, []string{"INVALID_REQUEST", "Invalid competitor name"}, fmt.Errorf("could not find competitor with given competitor name: %w", err))
				return
			}
			errorResponse(w, http.StatusInternalServerError, []string{"DATABASE_ERROR", "Something went wrong"}, fmt.Errorf("could not get competitor from the database: %w", err))
			return
		}

		if sportID != home.SportID {
			errorResponse(w, http.StatusBadRequest, []string{"INVALID_REQUEST", "Invalid data for given sport"}, nil)
			return
		}
		homeID = sql.NullInt32{
			Int32: home.ID,
			Valid: true,
		}
	}
	if reqData.AwayName != nil {
		away, err := cfg.Q.GetCompetitorByName(r.Context(), *reqData.AwayName)
		if err != nil {
			if err == sql.ErrNoRows {
				errorResponse(w, http.StatusBadRequest, []string{"INVALID_REQUEST", "Invalid competitor name"}, fmt.Errorf("could not find competitor with given competitor name: %w", err))
				return
			}
			errorResponse(w, http.StatusInternalServerError, []string{"DATABASE_ERROR", "Something went wrong"}, fmt.Errorf("could not get competitor from the database: %w", err))
			return
		}

		if sportID != away.SportID {
			errorResponse(w, http.StatusBadRequest, []string{"INVALID_REQUEST", "Invalid data for given sport"}, nil)
			return
		}
		awayID = sql.NullInt32{
			Int32: away.ID,
			Valid: true,
		}
	}

	stage, err := cfg.Q.GetStageByID(r.Context(), reqData.StageID)
	if err != nil {
		if err == sql.ErrNoRows {
			errorResponse(w, http.StatusBadRequest, []string{"INVALID_REQUEST", "Invalid competition stage"}, fmt.Errorf("could not find stage with given stage ID: %w", err))
			return
		}
		errorResponse(w, http.StatusInternalServerError, []string{"DATABASE_ERROR", "Something went wrong"}, fmt.Errorf("could not get stage from the database: %w", err))
		return
	}

	if comp, err := cfg.Q.GetCompByStageID(r.Context(), stage.ID); err != nil {
		errorResponse(w, http.StatusInternalServerError, []string{"DATABASE_ERROR", "Something went wrong"}, fmt.Errorf("could not get stage from the database: %w", err))
		return
	} else if comp.SportID != sportID {
		errorResponse(w, http.StatusBadRequest, []string{"INVALID_REQUEST", "Invalid data for given sport"}, nil)
		return
	}

	if reqData.HomeName != nil || reqData.AwayName != nil {
		constraintParams := database.GetEventByConstraintParams{
			VenueDate:        reqData.VenueDate,
			StageID:          stage.ID,
			HomeCompetitorID: homeID,
			AwayCompetitorID: awayID,
		}
		if err := cfg.Q.GetEventByConstraint(r.Context(), constraintParams); err == nil {
			errorResponse(w, http.StatusConflict, []string{"CONFLICTED_DATA", "Similar event already exists"}, nil)
			return
		}
	}

	tx, err := cfg.Db.BeginTx(r.Context(), nil)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, []string{"DATABASE_ERROR", "Something went wrong"}, fmt.Errorf("could not start database transaction: %w", err))
		return
	}
	defer tx.Rollback()

	qtx := cfg.Q.WithTx(tx)

	newEventParams := database.CreateEventParams{
		Status:           reqData.Status,
		VenueDate:        reqData.VenueDate,
		VenueTime:        reqData.VenueTime,
		VenueID:          venueID,
		HomeCompetitorID: homeID,
		AwayCompetitorID: awayID,
		StageID:          stage.ID,
		Details:          reqData.Details,
	}
	eventID, err := qtx.CreateEvent(r.Context(), newEventParams)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, []string{"DATABASE_ERROR", "Something went wrong"}, fmt.Errorf("could not create event in the database: %w", err))
	}

	if reqData.Result != nil {
		result := reqData.Result

		if result.Outcome != nil && !slices.Contains(availableOutcomes, *result.Outcome) {
			errorResponse(w, http.StatusBadRequest, []string{"INVALID_REQUEST", "Invalid result outcome"}, nil)
			return
		}
		if result.ForteitBy != nil && !slices.Contains(availableForfeiters, *result.ForteitBy) {
			errorResponse(w, http.StatusBadRequest, []string{"INVALID_REQUEST", "Invalid result forfeiter"}, nil)
			return
		}
		if (result.Outcome != nil && *result.Outcome == "forfeit") && (result.ForteitBy == nil && *result.ForteitBy == "") {
			errorResponse(w, http.StatusBadRequest, []string{"INVALID_REQUEST", "Invalid result outcome"}, nil)
			return
		}
		newResultParams := database.CreateResultParams{
			HomePoints: result.HomePoints,
			AwayPoints: result.AwayPoints,
			Outcome:    result.Outcome,
			ForfeitBy:  result.ForteitBy,
			EventID:    eventID,
			Details:    result.ResultDetails,
		}
		_, err := qtx.CreateResult(r.Context(), newResultParams)
		if err != nil {
			errorResponse(w, http.StatusInternalServerError, []string{"DATABASE_ERROR", "Something went wrong"}, fmt.Errorf("could not create result in the database: %w", err))
			return
		}
	}

	if err := tx.Commit(); err != nil {
		errorResponse(w, http.StatusInternalServerError, []string{"DATABASE_ERROR", "Something went wrong"}, fmt.Errorf("could not commit database transaction: %w", err))
		return
	}

	successResponse(w, http.StatusCreated, struct {
		EventID int32 `json:"event_id"`
	}{
		EventID: eventID,
	})
}

func (cfg ApiConfig) getEventHandler(w http.ResponseWriter, r *http.Request) {
	eventID := r.PathValue("eventID")
	if eventID == "" {
		errorResponse(w, http.StatusBadRequest, []string{"INVALID_REQUEST", "Invalid event ID"}, nil)
		return
	}
	intEventID, err := strconv.Atoi(eventID)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, []string{"INVALID_REQUEST", "Invalid event ID"}, fmt.Errorf("could not parse event ID from the URL path value: %w", err))
		return
	}
	event, err := cfg.Q.GetEventByID(r.Context(), int32(intEventID))
	if err != nil {
		if err == sql.ErrNoRows {
			errorResponse(w, http.StatusNotFound, []string{"NOT_FOUND", "Event not found"}, fmt.Errorf("could not find event with given event ID: %w", err))
			return
		}
		errorResponse(w, http.StatusInternalServerError, []string{"DATABASE_ERROR", "Something went wrong"}, fmt.Errorf("could not get event from the database: %w", err))
		return
	}

	venueDate, venueTime := formatDateTime(event.VenueDate, event.VenueTime)

	winner := computeWinner(event.HomePoints, event.AwayPoints, event.Status, event.Outcome, event.ForfeitBy)

	var venue *Venue
	if event.Country != nil || event.City != nil || event.PlaceName != nil {
		var countryPointer *string
		if event.Country != nil {
			eventCountry := strings.TrimSpace(*event.Country)
			countryPointer = &eventCountry
		}
		venue = &Venue{
			CountryCode: countryPointer,
			City:        event.City,
			Place:       event.PlaceName,
		}
	}

	eventRes := Event{
		Status:      event.Status,
		Date:        venueDate,
		Time:        venueTime,
		Venue:       venue,
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
	events, err := cfg.Q.GetEvents(r.Context(), getEventsParams)
	if err != nil {
		if err == sql.ErrNoRows {
			successResponse(w, http.StatusOK, []Event{})
			return
		}
		errorResponse(w, http.StatusInternalServerError, []string{"DATABASE_ERROR", "Something went wrong"}, fmt.Errorf("could not get events from the database: %w", err))
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
			Status:      ev.Status,
			Date:        venueDate,
			Time:        venueTime,
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
