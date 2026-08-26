-- 1. Restaurant Reviews
CREATE TABLE IF NOT EXISTS restaurants.reviews (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    restaurant_id UUID NOT NULL REFERENCES restaurants.restaurants(id) ON DELETE CASCADE,
    order_id UUID NOT NULL UNIQUE,
    customer_id UUID NOT NULL,
    score SMALLINT NOT NULL CHECK (score >= 1 AND score <= 5),
    comment TEXT,
    image_url TEXT,
    owner_reply TEXT,
    owner_reply_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_reviews_restaurant_id ON restaurants.reviews (restaurant_id);
CREATE INDEX IF NOT EXISTS idx_reviews_score ON restaurants.reviews (score);

-- 2. Item Reviews (Optional, for item-level ratings)
CREATE TABLE IF NOT EXISTS restaurants.review_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    review_id UUID NOT NULL REFERENCES restaurants.reviews(id) ON DELETE CASCADE,
    menu_item_id UUID NOT NULL REFERENCES restaurants.menu_items(id) ON DELETE CASCADE,
    score SMALLINT NOT NULL CHECK (score >= 1 AND score <= 5)
);

CREATE INDEX IF NOT EXISTS idx_review_items_menu_item ON restaurants.review_items (menu_item_id);
