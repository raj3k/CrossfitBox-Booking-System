CREATE TABLE IF NOT EXISTS workouts(
    id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    name text NOT NULL UNIQUE,
    mode text NOT NULL,
    time_cap integer NOT NULL,
    equipment text[],
    exercises text[] NOT NULL,
    trainer_tips text[],
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);

ALTER TABLE workouts ADD CONSTRAINT time_cap_check CHECK (time_cap >= 0);