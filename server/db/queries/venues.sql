-- name: AddVenue :one
INSERT INTO venues (name, _country, city, _sport_id)
VALUES (
  $1,
  $2,
  $3,
  $4
)
RETURNING id;

-- name: GetVenueByID :one
SELECT * FROM venues
WHERE id = $1;

-- name: GetVenuesForSport :many
SELECT 
  v.id,
  v.name

FROM venues v
LEFT JOIN sports s ON s.id = v._sport_id

WHERE s.name = $1

ORDER BY v.name;
