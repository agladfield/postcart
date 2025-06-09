-- name: SetInboundEmail :exec
INSERT INTO postmark_inbound_emails (id, received, email, from_name, subject, message, attachment)
VALUES (?, ?, ?, ?, ?, ?, ?)
ON CONFLICT (id) DO UPDATE
SET received = excluded.received,
    email = excluded.email,
    from_name = excluded.from_name,
    subject = excluded.subject,
    message = excluded.message,
    attachment = excluded.attachment;

-- name: DeleteSession :exec
DELETE FROM postmark_inbound_rules WHERE id = ?;

-- name: ResetSessions :exec
DELETE FROM postmark_inbound_rules;
