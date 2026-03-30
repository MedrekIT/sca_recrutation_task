-- name: AddCompetitor :one
INSERT INTO competitors (name, _country, _sport_id)
VALUES (
  $1,
  $2,
  $3
)
RETURNING id;

-- name: GetCompetitorByID :one
SELECT * FROM competitors
WHERE id = $1;

-- name: GetCompetitorByName :one
SELECT * FROM competitors
WHERE name = $1;

-- name: GetCompetitorsForSport :many
SELECT 
  c.name

FROM competitors c
LEFT JOIN sports s ON s.id = c._sport_id

WHERE s.name = $1

ORDER BY c.name;
