-- +goose Up
CREATE TABLE sports (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  name TEXT UNIQUE NOT NULL
);

-- +goose Down
DROP TABLE sports;
