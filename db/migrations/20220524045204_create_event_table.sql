-- migrate:up

CREATE TABLE event (
    id uuid PRIMARY KEY default gen_random_uuid(),
    event jsonb
);

-- migrate:down

DROP TABLE event;