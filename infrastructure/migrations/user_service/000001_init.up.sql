-- ============================================================
-- Schema: users
-- Tables: users, refresh_tokens, user_addresses
-- ============================================================

SET search_path TO users, public;

-- 1. users
CREATE TABLE IF NOT EXISTS users.users (
    id            UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    name          VARCHAR(100) NOT NULL,
    email         VARCHAR(150) NOT NULL,
    phone         VARCHAR(15),
    password_hash VARCHAR(255) NOT NULL,
    role          VARCHAR(20)  NOT NULL CHECK (role IN ('customer', 'restaurant_owner', 'admin')),
    avatar_url    TEXT,
    is_verified   BOOLEAN     NOT NULL DEFAULT FALSE,
    is_active     BOOLEAN     NOT NULL DEFAULT TRUE,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT users_email_unique UNIQUE (email),
    CONSTRAINT users_phone_unique UNIQUE (phone)
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users.users (email);
CREATE INDEX IF NOT EXISTS idx_users_role  ON users.users (role);

-- 2. refresh_tokens
CREATE TABLE IF NOT EXISTS users.refresh_tokens (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID        NOT NULL REFERENCES users.users (id) ON DELETE CASCADE,
    token_hash  VARCHAR(255) NOT NULL,
    expires_at  TIMESTAMPTZ NOT NULL,
    revoked_at  TIMESTAMPTZ,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT refresh_tokens_token_hash_unique UNIQUE (token_hash)
);

CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_id ON users.refresh_tokens (user_id);

-- 3. user_addresses
CREATE TABLE IF NOT EXISTS users.user_addresses (
    id          UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID         NOT NULL REFERENCES users.users (id) ON DELETE CASCADE,
    label       VARCHAR(50),
    street      TEXT         NOT NULL,
    district    VARCHAR(100),
    city        VARCHAR(100),
    lat         DECIMAL(10, 7),
    lng         DECIMAL(10, 7),
    is_default  BOOLEAN      NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_user_addresses_user_id ON users.user_addresses (user_id);
