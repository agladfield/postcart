-- -- name: UpdateSettings :exec
-- INSERT INTO user_config (id, updated, map_type, units_type, notifications_enabled, dark_mode, default_radius)
-- VALUES (?, ?, ?, ?, ?, ?, ?)
-- ON CONFLICT (id) DO UPDATE
-- SET updated = excluded.updated,
--     map_type = excluded.map_type,
--     units_type = excluded.units_type,
--     notifications_enabled = excluded.notifications_enabled,
--     dark_mode = excluded.dark_mode,
--     default_radius = excluded.default_radius;

-- name: CreateSender :exec
INSERT INTO senders (id, created, last_sent, email, sent, fails, delivered, blocked)
VALUES (?, ?, ?, ?, ?, ?, ?, ?);

-- name: GetSenderByID :one
SELECT * FROM senders WHERE id = ?;

-- name: GetSenderByEmail :one
SELECT * FROM senders WHERE email = ?;

-- name: CreateRecipient :exec
INSERT INTO recipients (id, created, email, received)
VALUES (?, ?, ?, ?);

-- name: GetRecipientByID :one
SELECT * FROM recipients WHERE id = ?;

-- name: GetRecipientByEmail :one
SELECT * FROM recipients WHERE email = ?;

-- name: IncrementRecipientByID :exec
UPDATE recipients
SET received = received + 1
WHERE id = ?;

-- name: IncrementRecipientByEmail :exec
UPDATE recipients
SET received = received + 1
WHERE email = ?;
