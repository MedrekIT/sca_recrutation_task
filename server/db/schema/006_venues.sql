-- +goose Up
CREATE TABLE venues (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  name TEXT NOT NULL,
  city TEXT,
  country_code CHAR(3)
);

-- +goose Down
DROP TABLE venues;
