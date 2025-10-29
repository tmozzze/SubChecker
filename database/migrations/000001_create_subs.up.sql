-- 000001_create_subs.up.sql

CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS subs (
    sub_id SERIAL PRIMARY KEY,
    service_name TEXT NOT NULL,
    price INT NOT NULL CHECK (price >= 0),
    user_id UUID NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NULL
);

CREATE INDEX IF NOT EXISTS idx_subs_user_id ON subs (user_id);
CREATE INDEX IF NOT EXISTS idx_subs_service_name ON subs (service_name);
CREATE INDEX IF NOT EXISTS idx_subs_start_date ON subs (start_date);
CREATE INDEX IF NOT EXISTS idx_subs_end_date ON subs (end_date);
