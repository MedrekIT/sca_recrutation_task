-- +goose Up
CREATE TABLE competitions (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  name TEXT NOT NULL,
  _sport_id INT NOT NULL REFERENCES sports(id) ON DELETE RESTRICT
);

-- +goose Down
DROP TABLE competitions;
