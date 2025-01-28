-- name: CreateJob :one
INSERT INTO jobs (name, handler, interval, status, last_triggered)
VALUES ($1, $2, $3, $4, $5)
    RETURNING *;

-- name: GetDueJobs :many
SELECT * FROM jobs
WHERE status = 'active' AND last_triggered + interval * interval '1 second' < NOW()
    FOR UPDATE SKIP LOCKED;

-- name: UpdateJobLastTriggered :exec
UPDATE jobs
SET last_triggered = $1
WHERE id = $2;

-- name: UpdateJobStatus :exec
UPDATE jobs
SET status = $1
WHERE id = $2;

-- name: GetNotSendedMessages :many
SELECT id, content, recipient_phone_number, status, message_received_id
FROM messages
WHERE status = 'pending'
ORDER BY created_at ASC
    LIMIT 2
FOR UPDATE SKIP LOCKED;


-- name: CreateMessage :one
INSERT INTO messages (id, content, recipient_phone_number, status, message_received_id)
VALUES ($1, $2, $3, $4, $5)
    RETURNING id, content, recipient_phone_number, status, message_received_id;

-- name: UpdateMessageStatus :exec
UPDATE messages
SET status = $1
WHERE id = $2;