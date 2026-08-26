-- ============================================================
-- Schema: notifications
-- Tables: notifications
-- ============================================================

CREATE TABLE IF NOT EXISTS notifications.notifications (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID        NOT NULL,       -- FK → users.users
    type        VARCHAR(50) NOT NULL,
    -- types: order_placed | order_confirmed | order_preparing |
    --        order_ready | order_completed | order_cancelled |
    --        payment_success | payment_failed | new_rating
    title       VARCHAR(200) NOT NULL,
    body        TEXT         NOT NULL,
    data        JSONB,                      -- {order_id, restaurant_id,...}
    is_read     BOOLEAN      NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_notifications_user_id  ON notifications.notifications (user_id);
CREATE INDEX IF NOT EXISTS idx_notifications_is_read  ON notifications.notifications (user_id, is_read);
CREATE INDEX IF NOT EXISTS idx_notifications_created  ON notifications.notifications (created_at DESC);
