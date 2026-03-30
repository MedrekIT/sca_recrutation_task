-- name: AddSport :one
INSERT INTO sports (name)
VALUES (
  $1
)
RETURNING *;

-- name: GetSportByID :one
SELECT * FROM sports
WHERE id = $1;

-- name: GetSportByName :one
SELECT * FROM sports
WHERE name = $1;

-- name: GetSports :many
SELECT name FROM sports;
