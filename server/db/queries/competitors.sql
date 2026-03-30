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
