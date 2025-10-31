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
