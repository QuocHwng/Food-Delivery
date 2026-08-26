-- ============================================================
-- Schema: audit
-- Tables: audit_logs
-- ============================================================

CREATE TABLE IF NOT EXISTS audit.audit_logs (
    id           UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id      UUID,                       -- NULL nếu là system action
    action       VARCHAR(100) NOT NULL,      -- 'order.cancel', 'menu.delete', 'user.login'...
    entity_type  VARCHAR(50),                -- 'order', 'user', 'menu_item'...
    entity_id    UUID,
    old_value    JSONB,                      -- giá trị trước khi thay đổi
    new_value    JSONB,                      -- giá trị sau khi thay đổi
    ip_address   VARCHAR(45),
    user_agent   TEXT,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_audit_user_id     ON audit.audit_logs (user_id);
CREATE INDEX IF NOT EXISTS idx_audit_action      ON audit.audit_logs (action);
CREATE INDEX IF NOT EXISTS idx_audit_entity      ON audit.audit_logs (entity_type, entity_id);
CREATE INDEX IF NOT EXISTS idx_audit_created_at  ON audit.audit_logs (created_at DESC);

-- Partition theo tháng sẽ thêm sau khi scale
