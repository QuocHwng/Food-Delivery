-- ============================================================
-- Schema: ratings
-- Tables: ratings, rating_replies
-- ============================================================

-- 1. ratings
CREATE TABLE IF NOT EXISTS ratings.ratings (
    id             UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id       UUID        NOT NULL,    -- FK → orders.orders
    customer_id    UUID        NOT NULL,    -- FK → users.users
    restaurant_id  UUID        NOT NULL,    -- FK → restaurants.restaurants
    menu_item_id   UUID,                    -- NULL = đánh giá chung, có = đánh giá món cụ thể
    food_score     SMALLINT    NOT NULL CHECK (food_score BETWEEN 1 AND 5),
    service_score  SMALLINT    NOT NULL CHECK (service_score BETWEEN 1 AND 5),
    comment        TEXT,
    images         JSONB,                   -- ["url1", "url2"]
    is_visible     BOOLEAN     NOT NULL DEFAULT TRUE,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT ratings_order_id_unique UNIQUE (order_id)  -- 1 đơn chỉ 1 đánh giá
);

CREATE INDEX IF NOT EXISTS idx_ratings_restaurant  ON ratings.ratings (restaurant_id);
CREATE INDEX IF NOT EXISTS idx_ratings_customer    ON ratings.ratings (customer_id);
CREATE INDEX IF NOT EXISTS idx_ratings_food_score  ON ratings.ratings (food_score);
CREATE INDEX IF NOT EXISTS idx_ratings_visible     ON ratings.ratings (is_visible, created_at DESC);

-- 2. rating_replies (owner trả lời review)
CREATE TABLE IF NOT EXISTS ratings.rating_replies (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    rating_id   UUID        NOT NULL REFERENCES ratings.ratings (id) ON DELETE CASCADE,
    owner_id    UUID        NOT NULL,    -- FK → users.users
    message     TEXT        NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT rating_replies_rating_unique UNIQUE (rating_id)  -- 1 đánh giá chỉ 1 reply
);
