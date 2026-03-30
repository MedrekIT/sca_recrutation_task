-- name: CreateEvent :one
INSERT INTO events (status, venue_date, venue_time, _venue_id, _home_competitor_id, _away_competitor_id, _stage_id, details)
VALUES (
  $1,
  $2,
  $3,
  $4,
  $5,
  $6,
  $7,
  $8
)
RETURNING id;

-- name: GetEventByConstraint :exec
SELECT * FROM events
WHERE venue_date = $1 AND _stage_id = $2 AND _home_competitor_id = $3 AND _away_competitor_id = $4;

-- name: GetEventByID :one
SELECT
  e.venue_date,
  e.venue_time,
  e.status,
  e.details AS event_details,

  v._country,
  v.city,
  v.name AS place_name,

  c.name AS competition,
  ed.season,
  st.name AS stage,

  hc._country AS home_country,
  hc.name AS home_competitor,
  ac._country AS away_country,
  ac.name AS away_competitor,

  r.home_points,
  r.away_points,
  r.outcome,
  r.forfeit_by,
  r.details AS result_details

FROM events e
JOIN stages st ON e._stage_id = st.id
JOIN editions ed ON st._edition_id = ed.id
JOIN competitions c ON ed._competition_id = c.id
LEFT JOIN venues v ON e._venue_id = v.id
LEFT JOIN competitors hc ON e._home_competitor_id = hc.id
LEFT JOIN competitors ac ON e._away_competitor_id = ac.id
LEFT JOIN results r ON r._event_id = e.id

WHERE e.id = $1;

-- name: GetEvents :many
SELECT
  e.id AS event_id,
  e.venue_date,
  e.venue_time,
  e.status,

  c.name AS competition,
  ed.season,
  st.name AS stage,

  hc.name AS home_competitor,
  ac.name AS away_competitor,

  r.home_points,
  r.away_points,
  r.outcome,
  r.forfeit_by

FROM events e
JOIN stages st ON e._stage_id = st.id
JOIN editions ed ON st._edition_id = ed.id
JOIN competitions c ON ed._competition_id = c.id
JOIN sports s ON c._sport_id = s.id
LEFT JOIN competitors hc ON e._home_competitor_id = hc.id
LEFT JOIN competitors ac ON e._away_competitor_id = ac.id
LEFT JOIN results r ON r._event_id = e.id

WHERE (sqlc.narg(date_filter)::date IS NULL OR e.venue_date = sqlc.narg(date_filter)::date) AND (sqlc.narg(sport_filter)::text IS NULL OR s.name = sqlc.narg(sport_filter)::text)

ORDER BY
  e.venue_date,
  e.venue_time,
  c.name;
