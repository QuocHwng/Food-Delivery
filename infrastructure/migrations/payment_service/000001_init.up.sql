-- ============================================================
-- Schema: payments
-- Tables: payments, payment_logs
-- VNPay Sandbox integration
-- ============================================================

-- 1. payments
CREATE TABLE IF NOT EXISTS payments.payments (
    id                    UUID          PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id              UUID          NOT NULL,        -- FK → orders.orders
    amount                DECIMAL(12,2) NOT NULL CHECK (amount > 0),
    method                VARCHAR(10)   NOT NULL DEFAULT 'vnpay'
                                        CHECK (method IN ('vnpay', 'cod')),
    status                VARCHAR(10)   NOT NULL DEFAULT 'pending'
                                        CHECK (status IN ('pending','paid','failed','refunded')),
    -- VNPay Sandbox fields
    vnpay_txn_ref         VARCHAR(100),    -- mã giao dịch mình tạo gửi lên VNPay
    vnpay_txn_no          VARCHAR(100),    -- mã giao dịch VNPay trả về
    vnpay_response_code   VARCHAR(10),     -- 00 = thành công
    vnpay_bank_code       VARCHAR(20),     -- VIETCOMBANK, NCB...
    vnpay_card_type       VARCHAR(20),     -- ATM, QR
    vnpay_pay_date        TIMESTAMPTZ,
    -- Refund
    refund_amount         DECIMAL(12,2),
    refund_at             TIMESTAMPTZ,
    refund_reason         TEXT,
    created_at            TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    updated_at            TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    CONSTRAINT payments_order_id_unique UNIQUE (order_id),
    CONSTRAINT payments_vnpay_txn_ref_unique UNIQUE (vnpay_txn_ref)
);

CREATE INDEX IF NOT EXISTS idx_payments_order_id    ON payments.payments (order_id);
CREATE INDEX IF NOT EXISTS idx_payments_status      ON payments.payments (status);
CREATE INDEX IF NOT EXISTS idx_payments_vnpay_ref   ON payments.payments (vnpay_txn_ref);

-- 2. payment_logs (log toàn bộ request/response VNPay để debug sandbox)
CREATE TABLE IF NOT EXISTS payments.payment_logs (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    payment_id  UUID        NOT NULL REFERENCES payments.payments (id) ON DELETE CASCADE,
    event       VARCHAR(50) NOT NULL,   -- 'initiate' | 'ipn' | 'return' | 'query' | 'refund'
    direction   VARCHAR(10) NOT NULL DEFAULT 'inbound'
                            CHECK (direction IN ('inbound', 'outbound')),
    payload     JSONB       NOT NULL,   -- raw request hoặc response
    ip_address  VARCHAR(45),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_payment_logs_payment ON payments.payment_logs (payment_id);
CREATE INDEX IF NOT EXISTS idx_payment_logs_event   ON payments.payment_logs (event);
