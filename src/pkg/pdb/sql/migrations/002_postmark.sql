-- +goose Up
CREATE TABLE IF NOT EXISTS postmark_inbound_emails (
    id TEXT PRIMARY KEY,
    received INTEGER NOT NULL,
    email TEXT NOT NULL,
    from_name TEXT NOT NULL,
    subject TEXT NOT NULL,
    message TEXT NOT NULL,
    attachment BLOB
) STRICT;

CREATE TABLE IF NOT EXISTS postmark_outbound_emails (
    id TEXT PRIMARY KEY,
    sent INTEGER NOT NULL,
    email TEXT NOT NULL,
    sender TEXT NOT NULL,
    status INTEGER NOT NULL
) STRICT;

CREATE TABLE IF NOT EXISTS postmark_bounced_emails (
    id TEXT PRIMARY KEY,
    bounced INTEGER NOT NULL,
    recipient TEXT NOT NULL,
    sender TEXT NOT NULL,
    description TEXT NOT NULL,
    details TEXT NOT NULL
) STRICT;

CREATE TABLE IF NOT EXISTS postmark_delivered_emails (
    id TEXT PRIMARY KEY,
    sender TEXT NOT NULL
) STRICT;

CREATE TABLE IF NOT EXISTS postmark_spam_complaints (
    id TEXT PRIMARY KEY
) STRICT;

CREATE TABLE IF NOT EXISTS postmark_templates (
    id TEXT PRIMARY KEY
) STRICT;

CREATE TABLE IF NOT EXISTS postmark_inbound_rules (
    id INTEGER PRIMARY KEY,
    created INTEGER NOT NULL,
    rule TEXT NOT NULL
) STRICT;

-- +goose Down
DROP TABLE IF EXISTS postmark_inbound_rules;
DROP TABLE IF EXISTS postmark_templates;
DROP TABLE IF EXISTS postmark_spam_complaints;
DROP TABLE IF EXISTS postmark_delivered_emails;
DROP TABLE IF EXISTS postmark_bounced_emails;
DROP TABLE IF EXISTS postmark_outbound_emails;
DROP TABLE IF EXISTS postmark_inbound_emails;
