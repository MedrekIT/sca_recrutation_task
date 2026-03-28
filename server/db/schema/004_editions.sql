-- +goose Up
CREATE TABLE editions (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  season TEXT NOT NULL,
  _competition_id INT NOT NULL REFERENCES competitions(id) ON DELETE CASCADE,
  UNIQUE (_competition_id, season)
);

-- +goose Down
DROP TABLE editions;
