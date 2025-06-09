-- +goose Up
CREATE TABLE IF NOT EXISTS queued_requests (
    id TEXT PRIMARY KEY,
    user TEXT NOT NULL,
    to_email TEXT NOT NULL,
    to_name TEXT NOT NULL,
    from_email TEXT NOT NULL,
    from_name TEXT NOT NULL,
    artwork_enum INTEGER NOT NULL,
    style_enum INTEGER NOT NULL,
    border_enum INTEGER NOT NULL,
    font_enum INTEGER NOT NULL,
    shape_enum INTEGER NOT NULL,
    country TEXT NOT NULL,
    subject TEXT NOT NULL,
    message TEXT NOT NULL,
    attachment BLOB
) STRICT;

CREATE TABLE IF NOT EXISTS completed_requests (
    id INTEGER PRIMARY KEY,
    user TEXT NOT NULL
) STRICT;

-- +goose Down
DROP TABLE IF EXISTS completed_requests;
DROP TABLE IF EXISTS active_requests;
