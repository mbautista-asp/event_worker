-- migrate:up

ALTER TABLE event ADD COLUMN created_at timestamptz default NOW();

-- migrate:down

ALTER TABLE event DROP COLUMN created_at;
