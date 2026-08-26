-- ============================================================
-- MASTER MIGRATION SCRIPT
-- Chạy toàn bộ: tạo schemas + tất cả bảng
-- psql -U hung -d food_delivery -f run_all_migrations.sql
-- Hoặc paste vào pgAdmin Query Tool
-- ============================================================

-- Extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Schemas
CREATE SCHEMA IF NOT EXISTS users;
CREATE SCHEMA IF NOT EXISTS restaurants;
CREATE SCHEMA IF NOT EXISTS orders;
CREATE SCHEMA IF NOT EXISTS payments;
CREATE SCHEMA IF NOT EXISTS ratings;
CREATE SCHEMA IF NOT EXISTS notifications;
CREATE SCHEMA IF NOT EXISTS audit;

-- ════════════════════════════════════════════
-- SCHEMA: users
-- ════════════════════════════════════════════

CREATE TABLE IF NOT EXISTS users.users (
    id            UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    name          VARCHAR(100) NOT NULL,
    email         VARCHAR(150) NOT NULL,
    phone         VARCHAR(15),
    password_hash VARCHAR(255) NOT NULL,
    role          VARCHAR(20)  NOT NULL CHECK (role IN ('customer','restaurant_owner','admin')),
    avatar_url    TEXT,
    is_verified   BOOLEAN      NOT NULL DEFAULT FALSE,
    is_active     BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    CONSTRAINT users_email_unique UNIQUE (email),
    CONSTRAINT users_phone_unique UNIQUE (phone)
);
CREATE INDEX IF NOT EXISTS idx_users_email ON users.users (email);
CREATE INDEX IF NOT EXISTS idx_users_role  ON users.users (role);

CREATE TABLE IF NOT EXISTS users.refresh_tokens (
    id          UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID         NOT NULL REFERENCES users.users (id) ON DELETE CASCADE,
    token_hash  VARCHAR(255) NOT NULL,
    expires_at  TIMESTAMPTZ  NOT NULL,
    revoked_at  TIMESTAMPTZ,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    CONSTRAINT refresh_tokens_hash_unique UNIQUE (token_hash)
);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user ON users.refresh_tokens (user_id);

CREATE TABLE IF NOT EXISTS users.user_addresses (
    id          UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID         NOT NULL REFERENCES users.users (id) ON DELETE CASCADE,
    label       VARCHAR(50),
    street      TEXT         NOT NULL,
    district    VARCHAR(100),
    city        VARCHAR(100),
    lat         DECIMAL(10,7),
    lng         DECIMAL(10,7),
    is_default  BOOLEAN      NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_user_addresses_user ON users.user_addresses (user_id);

-- ════════════════════════════════════════════
-- SCHEMA: restaurants
-- ════════════════════════════════════════════

CREATE TABLE IF NOT EXISTS restaurants.restaurant_categories (
    id         UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    name       VARCHAR(100) NOT NULL,
    slug       VARCHAR(100) NOT NULL,
    icon_url   TEXT,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    CONSTRAINT restaurant_categories_slug_unique UNIQUE (slug)
);

CREATE TABLE IF NOT EXISTS restaurants.restaurants (
    id               UUID          PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_id         UUID          NOT NULL,
    name             VARCHAR(200)  NOT NULL,
    description      TEXT,
    address          TEXT,
    lat              DECIMAL(10,7),
    lng              DECIMAL(10,7),
    phone            VARCHAR(15),
    logo_url         TEXT,
    banner_url       TEXT,
    avg_rating       DECIMAL(2,1)  NOT NULL DEFAULT 0,
    total_ratings    INT           NOT NULL DEFAULT 0,
    min_order_value  DECIMAL(12,2) NOT NULL DEFAULT 0,
    delivery_fee     DECIMAL(10,2) NOT NULL DEFAULT 0,
    is_open          BOOLEAN       NOT NULL DEFAULT TRUE,
    is_active        BOOLEAN       NOT NULL DEFAULT TRUE,
    created_at       TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_restaurants_owner    ON restaurants.restaurants (owner_id);
CREATE INDEX IF NOT EXISTS idx_restaurants_is_open  ON restaurants.restaurants (is_open);
CREATE INDEX IF NOT EXISTS idx_restaurants_rating   ON restaurants.restaurants (avg_rating DESC);

CREATE TABLE IF NOT EXISTS restaurants.restaurant_category_map (
    restaurant_id UUID NOT NULL REFERENCES restaurants.restaurants (id) ON DELETE CASCADE,
    category_id   UUID NOT NULL REFERENCES restaurants.restaurant_categories (id) ON DELETE CASCADE,
    PRIMARY KEY (restaurant_id, category_id)
);

CREATE TABLE IF NOT EXISTS restaurants.restaurant_operating_hours (
    id             UUID      PRIMARY KEY DEFAULT gen_random_uuid(),
    restaurant_id  UUID      NOT NULL REFERENCES restaurants.restaurants (id) ON DELETE CASCADE,
    day_of_week    SMALLINT  NOT NULL CHECK (day_of_week BETWEEN 0 AND 6),
    open_time      TIME,
    close_time     TIME,
    is_closed      BOOLEAN   NOT NULL DEFAULT FALSE,
    CONSTRAINT uq_operating_hours UNIQUE (restaurant_id, day_of_week)
);

CREATE TABLE IF NOT EXISTS restaurants.menu_categories (
    id             UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    restaurant_id  UUID         NOT NULL REFERENCES restaurants.restaurants (id) ON DELETE CASCADE,
    name           VARCHAR(100) NOT NULL,
    display_order  SMALLINT     NOT NULL DEFAULT 0,
    is_active      BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at     TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_menu_cat_restaurant ON restaurants.menu_categories (restaurant_id);

CREATE TABLE IF NOT EXISTS restaurants.menu_items (
    id             UUID          PRIMARY KEY DEFAULT gen_random_uuid(),
    restaurant_id  UUID          NOT NULL REFERENCES restaurants.restaurants (id) ON DELETE CASCADE,
    category_id    UUID          REFERENCES restaurants.menu_categories (id) ON DELETE SET NULL,
    name           VARCHAR(200)  NOT NULL,
    description    TEXT,
    price          DECIMAL(12,2) NOT NULL CHECK (price >= 0),
    image_url      TEXT,
    is_available   BOOLEAN       NOT NULL DEFAULT TRUE,
    display_order  SMALLINT      NOT NULL DEFAULT 0,
    avg_rating     DECIMAL(2,1)  NOT NULL DEFAULT 0,
    total_ratings  INT           NOT NULL DEFAULT 0,
    created_at     TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_menu_items_restaurant ON restaurants.menu_items (restaurant_id);
CREATE INDEX IF NOT EXISTS idx_menu_items_category   ON restaurants.menu_items (category_id);
CREATE INDEX IF NOT EXISTS idx_menu_items_available  ON restaurants.menu_items (is_available);

CREATE TABLE IF NOT EXISTS restaurants.menu_item_options (
    id            UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    menu_item_id  UUID         NOT NULL REFERENCES restaurants.menu_items (id) ON DELETE CASCADE,
    name          VARCHAR(100) NOT NULL,
    is_required   BOOLEAN      NOT NULL DEFAULT FALSE,
    max_select    SMALLINT     NOT NULL DEFAULT 1
);
CREATE INDEX IF NOT EXISTS idx_options_menu_item ON restaurants.menu_item_options (menu_item_id);

CREATE TABLE IF NOT EXISTS restaurants.menu_item_option_values (
    id          UUID          PRIMARY KEY DEFAULT gen_random_uuid(),
    option_id   UUID          NOT NULL REFERENCES restaurants.menu_item_options (id) ON DELETE CASCADE,
    label       VARCHAR(100)  NOT NULL,
    extra_price DECIMAL(10,2) NOT NULL DEFAULT 0
);
CREATE INDEX IF NOT EXISTS idx_option_values_option ON restaurants.menu_item_option_values (option_id);

-- Seed categories
INSERT INTO restaurants.restaurant_categories (name, slug) VALUES
    ('Cơm',         'com'),
    ('Bún - Phở',   'bun-pho'),
    ('Fast Food',   'fast-food'),
    ('Pizza',       'pizza'),
    ('Bánh mì',     'banh-mi'),
    ('Đồ uống',     'do-uong'),
    ('Tráng miệng', 'trang-mieng'),
    ('Lẩu - Nướng', 'lau-nuong'),
    ('Healthy',     'healthy'),
    ('Ăn vặt',      'an-vat')
ON CONFLICT (slug) DO NOTHING;

-- ════════════════════════════════════════════
-- SCHEMA: orders
-- ════════════════════════════════════════════

CREATE TABLE IF NOT EXISTS orders.coupons (
    id               UUID          PRIMARY KEY DEFAULT gen_random_uuid(),
    code             VARCHAR(50)   NOT NULL,
    type             VARCHAR(10)   NOT NULL CHECK (type IN ('percent','fixed')),
    value            DECIMAL(10,2) NOT NULL CHECK (value > 0),
    min_order_value  DECIMAL(12,2) NOT NULL DEFAULT 0,
    max_discount     DECIMAL(12,2),
    restaurant_id    UUID,
    max_uses         INT           NOT NULL DEFAULT -1,
    used_count       INT           NOT NULL DEFAULT 0,
    starts_at        TIMESTAMPTZ,
    expires_at       TIMESTAMPTZ,
    is_active        BOOLEAN       NOT NULL DEFAULT TRUE,
    created_at       TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    CONSTRAINT coupons_code_unique UNIQUE (code)
);
CREATE INDEX IF NOT EXISTS idx_coupons_code    ON orders.coupons (code);
CREATE INDEX IF NOT EXISTS idx_coupons_active  ON orders.coupons (is_active, expires_at);

CREATE TABLE IF NOT EXISTS orders.orders (
    id               UUID          PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_id      UUID          NOT NULL,
    restaurant_id    UUID          NOT NULL,
    address_id       UUID,
    status           VARCHAR(20)   NOT NULL DEFAULT 'pending'
                                   CHECK (status IN ('pending','confirmed','preparing','ready','completed','cancelled')),
    subtotal         DECIMAL(12,2) NOT NULL CHECK (subtotal >= 0),
    discount_amount  DECIMAL(12,2) NOT NULL DEFAULT 0,
    delivery_fee     DECIMAL(10,2) NOT NULL DEFAULT 0,
    total            DECIMAL(12,2) NOT NULL CHECK (total >= 0),
    note             TEXT,
    coupon_id        UUID          REFERENCES orders.coupons (id) ON DELETE SET NULL,
    cancelled_by     UUID,
    cancel_reason    TEXT,
    created_at       TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_orders_customer    ON orders.orders (customer_id);
CREATE INDEX IF NOT EXISTS idx_orders_restaurant  ON orders.orders (restaurant_id);
CREATE INDEX IF NOT EXISTS idx_orders_status      ON orders.orders (status);
CREATE INDEX IF NOT EXISTS idx_orders_created     ON orders.orders (created_at DESC);

CREATE TABLE IF NOT EXISTS orders.order_items (
    id               UUID          PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id         UUID          NOT NULL REFERENCES orders.orders (id) ON DELETE CASCADE,
    menu_item_id     UUID          NOT NULL,
    name_snapshot    VARCHAR(200)  NOT NULL,
    price_snapshot   DECIMAL(12,2) NOT NULL,
    quantity         SMALLINT      NOT NULL CHECK (quantity > 0),
    options_snapshot JSONB
);
CREATE INDEX IF NOT EXISTS idx_order_items_order ON orders.order_items (order_id);

CREATE TABLE IF NOT EXISTS orders.order_status_history (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id    UUID        NOT NULL REFERENCES orders.orders (id) ON DELETE CASCADE,
    from_status VARCHAR(20),
    to_status   VARCHAR(20) NOT NULL,
    changed_by  UUID        NOT NULL,
    note        TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_status_hist_order ON orders.order_status_history (order_id);

CREATE TABLE IF NOT EXISTS orders.coupon_usages (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    coupon_id   UUID        NOT NULL REFERENCES orders.coupons (id) ON DELETE CASCADE,
    user_id     UUID        NOT NULL,
    order_id    UUID        NOT NULL REFERENCES orders.orders (id) ON DELETE CASCADE,
    used_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT coupon_usages_unique UNIQUE (coupon_id, user_id, order_id)
);
CREATE INDEX IF NOT EXISTS idx_coupon_usages_coupon ON orders.coupon_usages (coupon_id);
CREATE INDEX IF NOT EXISTS idx_coupon_usages_user   ON orders.coupon_usages (user_id);

-- ════════════════════════════════════════════
-- SCHEMA: payments
-- ════════════════════════════════════════════

CREATE TABLE IF NOT EXISTS payments.payments (
    id                    UUID          PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id              UUID          NOT NULL,
    amount                DECIMAL(12,2) NOT NULL CHECK (amount > 0),
    method                VARCHAR(10)   NOT NULL DEFAULT 'vnpay'
                                        CHECK (method IN ('vnpay','cod')),
    status                VARCHAR(10)   NOT NULL DEFAULT 'pending'
                                        CHECK (status IN ('pending','paid','failed','refunded')),
    vnpay_txn_ref         VARCHAR(100),
    vnpay_txn_no          VARCHAR(100),
    vnpay_response_code   VARCHAR(10),
    vnpay_bank_code       VARCHAR(20),
    vnpay_card_type       VARCHAR(20),
    vnpay_pay_date        TIMESTAMPTZ,
    refund_amount         DECIMAL(12,2),
    refund_at             TIMESTAMPTZ,
    refund_reason         TEXT,
    created_at            TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    updated_at            TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    CONSTRAINT payments_order_unique    UNIQUE (order_id),
    CONSTRAINT payments_vnpay_ref_unique UNIQUE (vnpay_txn_ref)
);
CREATE INDEX IF NOT EXISTS idx_payments_order   ON payments.payments (order_id);
CREATE INDEX IF NOT EXISTS idx_payments_status  ON payments.payments (status);
CREATE INDEX IF NOT EXISTS idx_payments_ref     ON payments.payments (vnpay_txn_ref);

CREATE TABLE IF NOT EXISTS payments.payment_logs (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    payment_id  UUID        NOT NULL REFERENCES payments.payments (id) ON DELETE CASCADE,
    event       VARCHAR(50) NOT NULL,
    direction   VARCHAR(10) NOT NULL DEFAULT 'inbound'
                            CHECK (direction IN ('inbound','outbound')),
    payload     JSONB       NOT NULL,
    ip_address  VARCHAR(45),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_pay_logs_payment ON payments.payment_logs (payment_id);
CREATE INDEX IF NOT EXISTS idx_pay_logs_event   ON payments.payment_logs (event);

-- ════════════════════════════════════════════
-- SCHEMA: ratings
-- ════════════════════════════════════════════

CREATE TABLE IF NOT EXISTS ratings.ratings (
    id             UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id       UUID        NOT NULL,
    customer_id    UUID        NOT NULL,
    restaurant_id  UUID        NOT NULL,
    menu_item_id   UUID,
    food_score     SMALLINT    NOT NULL CHECK (food_score BETWEEN 1 AND 5),
    service_score  SMALLINT    NOT NULL CHECK (service_score BETWEEN 1 AND 5),
    comment        TEXT,
    images         JSONB,
    is_visible     BOOLEAN     NOT NULL DEFAULT TRUE,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT ratings_order_unique UNIQUE (order_id)
);
CREATE INDEX IF NOT EXISTS idx_ratings_restaurant ON ratings.ratings (restaurant_id);
CREATE INDEX IF NOT EXISTS idx_ratings_customer   ON ratings.ratings (customer_id);
CREATE INDEX IF NOT EXISTS idx_ratings_score      ON ratings.ratings (food_score, service_score);

CREATE TABLE IF NOT EXISTS ratings.rating_replies (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    rating_id   UUID        NOT NULL REFERENCES ratings.ratings (id) ON DELETE CASCADE,
    owner_id    UUID        NOT NULL,
    message     TEXT        NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT rating_replies_rating_unique UNIQUE (rating_id)
);

-- ════════════════════════════════════════════
-- SCHEMA: notifications
-- ════════════════════════════════════════════

CREATE TABLE IF NOT EXISTS notifications.notifications (
    id          UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID         NOT NULL,
    type        VARCHAR(50)  NOT NULL,
    title       VARCHAR(200) NOT NULL,
    body        TEXT         NOT NULL,
    data        JSONB,
    is_read     BOOLEAN      NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_notif_user    ON notifications.notifications (user_id);
CREATE INDEX IF NOT EXISTS idx_notif_unread  ON notifications.notifications (user_id, is_read);
CREATE INDEX IF NOT EXISTS idx_notif_created ON notifications.notifications (created_at DESC);

-- ════════════════════════════════════════════
-- SCHEMA: audit
-- ════════════════════════════════════════════

CREATE TABLE IF NOT EXISTS audit.audit_logs (
    id           UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id      UUID,
    action       VARCHAR(100) NOT NULL,
    entity_type  VARCHAR(50),
    entity_id    UUID,
    old_value    JSONB,
    new_value    JSONB,
    ip_address   VARCHAR(45),
    user_agent   TEXT,
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_audit_user    ON audit.audit_logs (user_id);
CREATE INDEX IF NOT EXISTS idx_audit_action  ON audit.audit_logs (action);
CREATE INDEX IF NOT EXISTS idx_audit_entity  ON audit.audit_logs (entity_type, entity_id);
CREATE INDEX IF NOT EXISTS idx_audit_time    ON audit.audit_logs (created_at DESC);

-- ════════════════════════════════════════════
-- DONE
-- ════════════════════════════════════════════
SELECT 'food_delivery database initialized successfully! 🎉' AS status;
