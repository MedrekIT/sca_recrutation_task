-- name: GetCompEditionsForSport :many
SELECT 
  c.name,
  e.id,
  e.season

FROM competitions c
JOIN editions e ON c.id = e._competition_id
LEFT JOIN sports s ON s.id = c._sport_id

WHERE s.name = $1

ORDER BY c.name;

-- name: GetCompByStageID :one
SELECT c.* FROM competitions c
JOIN editions e ON c.id = e._competition_id
JOIN stages s ON e.id = s._edition_id
WHERE s.id = $1;

-- name: GetStageByID :one
SELECT * FROM stages
WHERE id = $1;

-- name: GetStagesForCompEdition :many
SELECT 
  id,
  name
FROM stages
WHERE _edition_id = $1
ORDER BY name;
