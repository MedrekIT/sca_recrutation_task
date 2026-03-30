-- name: AddCountry :one
INSERT INTO countries (country_code, name)
VALUES (
  $1,
  $2
)
RETURNING *;

-- name: GetCountryByCode :one
SELECT * FROM countries
WHERE country_code = $1;
