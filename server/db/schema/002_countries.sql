-- +goose Up
CREATE TABLE countries (
  country_code CHAR(3) PRIMARY KEY,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  name TEXT UNIQUE NOT NULL
);

-- +goose Down
DROP TABLE countries;
