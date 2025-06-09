-- name: GetQueuedRequest :one
SELECT * FROM queued_requests WHERE id = 1;

-- name: GetQueuedRequests :many
SELECT * FROM queued_requests;

-- name: SetQueuedRequest :exec
INSERT INTO queued_requests (id, user, to_email, to_name, from_email, from_name, artwork_enum, style_enum, border_enum, font_enum, shape_enum, country, subject, message, attachment)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
ON CONFLICT (id) DO UPDATE
SET user = excluded.user,
    to_email = excluded.to_email,
    to_name = excluded.to_name,
    from_email = excluded.from_email,
    from_name = excluded.from_name,
    artwork_enum = excluded.artwork_enum,
    style_enum = excluded.style_enum,
    border_enum = excluded.border_enum,
    font_enum = excluded.font_enum,
    shape_enum = excluded.shape_enum,
    country = excluded.country,
    subject = excluded.subject,
    message = excluded.message,
    attachment = excluded.attachment;
