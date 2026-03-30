-- name: CreateResult :exec
INSERT INTO results (home_points, away_points, outcome, forfeit_by, _event_id, details)
VALUES (
  $1,
  $2,
  $3,
  $4,
  $5,
  $6
);
