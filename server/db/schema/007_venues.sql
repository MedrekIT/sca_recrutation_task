-- +goose Up
CREATE TABLE venues (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  name TEXT NOT NULL,
  city TEXT,
  _country CHAR(3) REFERENCES countries(country_code) ON DELETE SET NULL,
  _sport_id INT NOT NULL REFERENCES sports(id) ON DELETE RESTRICT
);

-- +goose Down
DROP TABLE venues;
