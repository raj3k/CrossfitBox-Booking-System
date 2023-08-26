CREATE INDEX IF NOT EXISTS workouts_name_idx ON workouts USING GIN (to_tsvector('simple', name));
CREATE INDEX IF NOT EXISTS workouts_mode_idx ON workouts USING GIN (to_tsvector('simple', mode));
CREATE INDEX IF NOT EXISTS workouts_equipment_idx ON workouts USING GIN (equipment);