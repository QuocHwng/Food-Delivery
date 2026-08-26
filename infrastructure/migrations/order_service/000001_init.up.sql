-- ============================================================
-- Schema: orders
-- Tables: coupons, orders, order_items,
--         order_status_history, coupon_usages
-- ============================================================

-- 1. coupons (tạo trước vì orders FK→coupons)
CREATE TABLE IF NOT EXISTS orders.coupons (
    id               UUID          PRIMARY KEY DEFAULT gen_random_uuid(),
    code             VARCHAR(50)   NOT NULL,
    type             VARCHAR(10)   NOT NULL CHECK (type IN ('percent', 'fixed')),
    value            DECIMAL(10,2) NOT NULL CHECK (value > 0),
    min_order_value  DECIMAL(12,2) NOT NULL DEFAULT 0,
    max_discount     DECIMAL(12,2),                  -- chặn trên cho percent
    restaurant_id    UUID,                            -- NULL = toàn hệ thống
    max_uses         INT           NOT NULL DEFAULT -1,  -- -1 = unlimited
    used_count       INT           NOT NULL DEFAULT 0,
    starts_at        TIMESTAMPTZ,
    expires_at       TIMESTAMPTZ,
    is_active        BOOLEAN       NOT NULL DEFAULT TRUE,
    created_at       TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    CONSTRAINT coupons_code_unique UNIQUE (code)
);

CREATE INDEX IF NOT EXISTS idx_coupons_code        ON orders.coupons (code);
CREATE INDEX IF NOT EXISTS idx_coupons_restaurant  ON orders.coupons (restaurant_id);
CREATE INDEX IF NOT EXISTS idx_coupons_active      ON orders.coupons (is_active, expires_at);

-- 2. orders
CREATE TABLE IF NOT EXISTS orders.orders (
    id               UUID          PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_id      UUID          NOT NULL,          -- FK → users.users
    restaurant_id    UUID          NOT NULL,          -- FK → restaurants.restaurants
    address_id       UUID,                            -- FK → users.user_addresses
    status           VARCHAR(20)   NOT NULL DEFAULT 'pending'
                                   CHECK (status IN ('pending','confirmed','preparing','ready','completed','cancelled')),
    subtotal         DECIMAL(12,2) NOT NULL CHECK (subtotal >= 0),
    discount_amount  DECIMAL(12,2) NOT NULL DEFAULT 0,
    delivery_fee     DECIMAL(10,2) NOT NULL DEFAULT 0,
    total            DECIMAL(12,2) NOT NULL CHECK (total >= 0),
    note             TEXT,
    coupon_id        UUID          REFERENCES orders.coupons (id) ON DELETE SET NULL,
    cancelled_by     UUID,                            -- user_id người huỷ
    cancel_reason    TEXT,
    created_at       TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_orders_customer_id   ON orders.orders (customer_id);
CREATE INDEX IF NOT EXISTS idx_orders_restaurant_id ON orders.orders (restaurant_id);
CREATE INDEX IF NOT EXISTS idx_orders_status        ON orders.orders (status);
CREATE INDEX IF NOT EXISTS idx_orders_created_at    ON orders.orders (created_at DESC);

-- 3. order_items (snapshot tại thời điểm đặt)
CREATE TABLE IF NOT EXISTS orders.order_items (
    id               UUID          PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id         UUID          NOT NULL REFERENCES orders.orders (id) ON DELETE CASCADE,
    menu_item_id     UUID          NOT NULL,          -- FK → restaurants.menu_items (snapshot)
    name_snapshot    VARCHAR(200)  NOT NULL,          -- tên món tại lúc đặt
    price_snapshot   DECIMAL(12,2) NOT NULL,          -- giá tại lúc đặt
    quantity         SMALLINT      NOT NULL CHECK (quantity > 0),
    options_snapshot JSONB                            -- tuỳ chọn đã chọn [{name, label, extra_price}]
);

CREATE INDEX IF NOT EXISTS idx_order_items_order_id ON orders.order_items (order_id);

-- 4. order_status_history (log mọi lần đổi trạng thái)
CREATE TABLE IF NOT EXISTS orders.order_status_history (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id    UUID        NOT NULL REFERENCES orders.orders (id) ON DELETE CASCADE,
    from_status VARCHAR(20),
    to_status   VARCHAR(20) NOT NULL,
    changed_by  UUID        NOT NULL,                -- user_id thực hiện đổi
    note        TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_status_history_order ON orders.order_status_history (order_id);

-- 5. coupon_usages (tracking ai đã dùng coupon nào)
CREATE TABLE IF NOT EXISTS orders.coupon_usages (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    coupon_id   UUID        NOT NULL REFERENCES orders.coupons (id) ON DELETE CASCADE,
    user_id     UUID        NOT NULL,                -- FK → users.users
    order_id    UUID        NOT NULL REFERENCES orders.orders (id) ON DELETE CASCADE,
    used_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT coupon_usages_user_order_unique UNIQUE (coupon_id, user_id, order_id)
);

CREATE INDEX IF NOT EXISTS idx_coupon_usages_coupon ON orders.coupon_usages (coupon_id);
CREATE INDEX IF NOT EXISTS idx_coupon_usages_user   ON orders.coupon_usages (user_id);
