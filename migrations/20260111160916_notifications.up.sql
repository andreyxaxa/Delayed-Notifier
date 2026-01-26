CREATE TABLE IF NOT EXISTS notifications
(
    uid         UUID        PRIMARY KEY,
    payload     JSONB       NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    send_at     TIMESTAMPTZ NOT NULL,
    status      TEXT        NOT NULL,
    retry_count INT         NOT NULL DEFAULT 0
);

CREATE INDEX IF NOT EXISTS idx_pending_notifications
ON notifications (send_at ASC)
WHERE status = 'pending';
