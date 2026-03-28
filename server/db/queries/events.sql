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
