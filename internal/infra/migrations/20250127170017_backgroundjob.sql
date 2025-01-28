-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE
EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE jobs
(
    id             UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    name           TEXT NOT NULL,
    handler        TEXT NOT NULL,
    interval       INT  NOT NULL,
    status         TEXT NOT NULL,
    last_triggered TIMESTAMP
);
CREATE INDEX idx_jobs_status ON jobs (status);
CREATE INDEX idx_jobs_last_triggered ON jobs (last_triggered);

INSERT INTO jobs (id,name, handler, interval, status, last_triggered)
VALUES (uuid_generate_v4(),'Job 1', 'message_publish_handler', 120, 'active', '2023-10-01 12:00:00');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_jobs_status;
DROP INDEX IF EXISTS idx_jobs_last_triggered;
SELECT 'down SQL query';
DROP TABLE jobs;

DROP
EXTENSION IF EXISTS "uuid-ossp";
-- +goose StatementEnd
