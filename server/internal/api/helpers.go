package api

import (
	"database/sql"
	"fmt"
	"time"
)

func formatDateTime(dbDate time.Time, dbTime *time.Time) (string, *string) {
	venueYear, venueMon, venueDay := dbDate.Date()
	venueDate := fmt.Sprintf("%d/%02d/%02d", venueYear, venueMon, venueDay)
	if dbTime == nil {
		return venueDate, nil
	}

	venueHour, venueMin, _ := dbTime.Clock()
	venueTime := fmt.Sprintf("%02d:%02d", venueHour, venueMin)

	return venueDate, &venueTime
}

func computeWinner(homePoints, awayPoints sql.NullInt16, status string, outcome, forfeitBy *string) string {
	if !homePoints.Valid || !awayPoints.Valid || status == "live" {
		return ""
	}

	if outcome == nil {
		if homePoints.Int16 > awayPoints.Int16 {
			return "home"
		} else if homePoints.Int16 < awayPoints.Int16 {
			return "away"
		}
	} else if *outcome != "" {
		status = "finished"
		if *outcome == "forfeit" && forfeitBy != nil {
			if *forfeitBy == "home" {
				return "away"
			}
			if *forfeitBy == "away" {
				return "home"
			}
		} else {
			return ""
		}
	}

	return ""
}
