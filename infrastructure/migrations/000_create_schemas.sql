-- ============================================================
-- FILE: infrastructure/migrations/000_create_database.sql
-- Tạo database + schemas + extensions
-- Chạy với user có quyền CREATEDB (hung)
-- ============================================================

-- Tạo database (bỏ nếu đã có)
-- Chạy riêng nếu cần: CREATE DATABASE food_delivery;

-- Bật extensions cần thiết
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Tạo các schemas (mỗi service 1 schema)
CREATE SCHEMA IF NOT EXISTS users;
CREATE SCHEMA IF NOT EXISTS restaurants;
CREATE SCHEMA IF NOT EXISTS orders;
CREATE SCHEMA IF NOT EXISTS payments;
CREATE SCHEMA IF NOT EXISTS ratings;
CREATE SCHEMA IF NOT EXISTS notifications;
CREATE SCHEMA IF NOT EXISTS audit;

-- Grant quyền cho user hung
GRANT ALL PRIVILEGES ON SCHEMA users         TO hung;
GRANT ALL PRIVILEGES ON SCHEMA restaurants   TO hung;
GRANT ALL PRIVILEGES ON SCHEMA orders        TO hung;
GRANT ALL PRIVILEGES ON SCHEMA payments      TO hung;
GRANT ALL PRIVILEGES ON SCHEMA ratings       TO hung;
GRANT ALL PRIVILEGES ON SCHEMA notifications TO hung;
GRANT ALL PRIVILEGES ON SCHEMA audit         TO hung;
