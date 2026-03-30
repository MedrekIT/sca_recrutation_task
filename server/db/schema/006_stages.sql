-- +goose Up
CREATE TABLE stages (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  name TEXT NOT NULL,
  _edition_id INT NOT NULL REFERENCES editions(id) ON DELETE CASCADE,
  UNIQUE (name, _edition_id)
);

-- +goose Down
DROP TABLE stages;
