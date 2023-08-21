CREATE TABLE IF NOT EXISTS workouts(
    id bigserial NOT NULL PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    title text NOT NULL,
    mode text NOT NULL,
    time_cap integer NOT NULL,
    equipment text[],
    exercises text[] NOT NULL,
    trainer_tips text[]
);

ALTER TABLE workouts ADD CONSTRAINT time_cap_check CHECK (time_cap >= 0);

CREATE INDEX IF NOT EXISTS workouts_id_title_mode_indx ON workouts (id, title, mode);