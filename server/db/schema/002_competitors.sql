-- +goose Up
CREATE TABLE competitors (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  name TEXT NOT NULL,
  country_code CHAR(3),
  type TEXT NOT NULL CHECK (type IN ('team', 'individual')),
  _sport_id INT NOT NULL REFERENCES sports(id) ON DELETE RESTRICT
);

-- +goose Down
DROP TABLE competitors;
