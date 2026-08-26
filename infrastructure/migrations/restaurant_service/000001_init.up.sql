-- ============================================================
-- Schema: restaurants
-- Tables: restaurants, restaurant_categories, restaurant_category_map,
--         restaurant_operating_hours, menu_categories,
--         menu_items, menu_item_options, menu_item_option_values
-- ============================================================

-- 1. restaurant_categories (danh mục: Fast food, Cơm, Bún phở...)
CREATE TABLE IF NOT EXISTS restaurants.restaurant_categories (
    id         UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    name       VARCHAR(100) NOT NULL,
    slug       VARCHAR(100) NOT NULL,
    icon_url   TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT restaurant_categories_slug_unique UNIQUE (slug)
);

-- 2. restaurants
CREATE TABLE IF NOT EXISTS restaurants.restaurants (
    id               UUID          PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_id         UUID          NOT NULL,              -- FK → users.users (cross-schema)
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

CREATE INDEX IF NOT EXISTS idx_restaurants_owner_id   ON restaurants.restaurants (owner_id);
CREATE INDEX IF NOT EXISTS idx_restaurants_is_open    ON restaurants.restaurants (is_open);
CREATE INDEX IF NOT EXISTS idx_restaurants_avg_rating ON restaurants.restaurants (avg_rating DESC);

-- 3. restaurant_category_map (many-to-many)
CREATE TABLE IF NOT EXISTS restaurants.restaurant_category_map (
    restaurant_id UUID NOT NULL REFERENCES restaurants.restaurants (id) ON DELETE CASCADE,
    category_id   UUID NOT NULL REFERENCES restaurants.restaurant_categories (id) ON DELETE CASCADE,
    PRIMARY KEY (restaurant_id, category_id)
);

-- 4. restaurant_operating_hours
CREATE TABLE IF NOT EXISTS restaurants.restaurant_operating_hours (
    id             UUID       PRIMARY KEY DEFAULT gen_random_uuid(),
    restaurant_id  UUID       NOT NULL REFERENCES restaurants.restaurants (id) ON DELETE CASCADE,
    day_of_week    SMALLINT   NOT NULL CHECK (day_of_week BETWEEN 0 AND 6), -- 0=CN,1=T2...6=T7
    open_time      TIME,
    close_time     TIME,
    is_closed      BOOLEAN    NOT NULL DEFAULT FALSE,
    CONSTRAINT uq_operating_hours UNIQUE (restaurant_id, day_of_week)
);

-- 5. menu_categories (danh mục trong menu: Khai vị, Món chính...)
CREATE TABLE IF NOT EXISTS restaurants.menu_categories (
    id             UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    restaurant_id  UUID         NOT NULL REFERENCES restaurants.restaurants (id) ON DELETE CASCADE,
    name           VARCHAR(100) NOT NULL,
    display_order  SMALLINT     NOT NULL DEFAULT 0,
    is_active      BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at     TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_menu_categories_restaurant ON restaurants.menu_categories (restaurant_id);

-- 6. menu_items
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

-- 7. menu_item_options (nhóm tuỳ chọn: Kích cỡ, Topping, Đá/Nóng...)
CREATE TABLE IF NOT EXISTS restaurants.menu_item_options (
    id            UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    menu_item_id  UUID         NOT NULL REFERENCES restaurants.menu_items (id) ON DELETE CASCADE,
    name          VARCHAR(100) NOT NULL,
    is_required   BOOLEAN      NOT NULL DEFAULT FALSE,
    max_select    SMALLINT     NOT NULL DEFAULT 1
);

CREATE INDEX IF NOT EXISTS idx_options_menu_item ON restaurants.menu_item_options (menu_item_id);

-- 8. menu_item_option_values (Nhỏ/Vừa/Lớn; Trân châu +5k...)
CREATE TABLE IF NOT EXISTS restaurants.menu_item_option_values (
    id          UUID          PRIMARY KEY DEFAULT gen_random_uuid(),
    option_id   UUID          NOT NULL REFERENCES restaurants.menu_item_options (id) ON DELETE CASCADE,
    label       VARCHAR(100)  NOT NULL,
    extra_price DECIMAL(10,2) NOT NULL DEFAULT 0
);

CREATE INDEX IF NOT EXISTS idx_option_values_option ON restaurants.menu_item_option_values (option_id);

-- Seed restaurant categories mặc định
INSERT INTO restaurants.restaurant_categories (name, slug) VALUES
    ('Cơm',           'com'),
    ('Bún - Phở',     'bun-pho'),
    ('Fast Food',     'fast-food'),
    ('Pizza',         'pizza'),
    ('Bánh mì',       'banh-mi'),
    ('Đồ uống',       'do-uong'),
    ('Tráng miệng',   'trang-mieng'),
    ('Lẩu - Nướng',   'lau-nuong'),
    ('Healthy',       'healthy'),
    ('Ăn vặt',        'an-vat')
ON CONFLICT (slug) DO NOTHING;
