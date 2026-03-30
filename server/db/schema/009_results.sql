-- +goose Up
CREATE TABLE results (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  home_points SMALLINT NOT NULL DEFAULT 0 CHECK (home_points >= 0),
  away_points SMALLINT NOT NULL DEFAULT 0 CHECK (away_points >= 0),
  outcome TEXT CHECK (outcome IN ('no_contest', 'forfeit')) DEFAULT NULL,
  forfeit_by TEXT CHECK (forfeit_by IN ('home', 'away')) DEFAULT NULL,
  _event_id INT UNIQUE NOT NULL REFERENCES events(id) ON DELETE CASCADE,
  details TEXT
);

-- +goose Down
DROP TABLE results;
