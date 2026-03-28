-- +goose Up
CREATE TABLE events (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  status TEXT NOT NULL CHECK (status IN ('scheduled', 'finished', 'live', 'cancelled', 'postponed')),
  _venue_id INT REFERENCES venues(id) ON DELETE SET NULL,
  venue_time TIME,
  venue_date DATE NOT NULL,
  _home_competitor_id INT REFERENCES competitors(id) ON DELETE RESTRICT,
  _away_competitor_id INT REFERENCES competitors(id) ON DELETE RESTRICT,
  _stage_id INT NOT NULL REFERENCES stages(id) ON DELETE RESTRICT,
  group_name TEXT,
  details TEXT,
  CONSTRAINT different_competitors CHECK (_home_competitor_id <> _away_competitor_id),
  UNIQUE (_home_competitor_id, _away_competitor_id, venue_date, _stage_id)
);

-- +goose Down
DROP TABLE events;
