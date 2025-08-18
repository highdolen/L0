-- migrations/01_create_tables.sql

CREATE TABLE delivery (
    id BIGSERIAL PRIMARY KEY,
    name TEXT,
    phone TEXT,
    zip TEXT,
    city TEXT,
    address TEXT,
    region TEXT,
    email TEXT
);

CREATE TABLE payment (
    id BIGSERIAL PRIMARY KEY,
    transaction TEXT,
    request_id TEXT,
    currency TEXT,
    provider TEXT,
    amount BIGINT,
    payment_dt BIGINT,
    bank TEXT,
    delivery_cost BIGINT,
    goods_total BIGINT,
    custom_fee BIGINT
);

CREATE TABLE orders (
    order_uid TEXT PRIMARY KEY,
    track_number TEXT,
    entry TEXT,
    delivery_id BIGINT REFERENCES deliveries(id) ON DELETE SET NULL,
    payment_id BIGINT REFERENCES payments(id) ON DELETE SET NULL,
    locale TEXT,
    internal_signature TEXT,
    customer_id TEXT,
    delivery_service TEXT,
    shardkey TEXT,
    sm_id INTEGER,
    date_created TIMESTAMPTZ,
    oof_shard TEXT
);

CREATE TABLE items (
    id BIGSERIAL PRIMARY KEY,
    chrt_id BIGINT,
    track_number TEXT,
    price BIGINT,
    rid TEXT,
    name TEXT,
    sale INTEGER,
    size TEXT,
    total_price BIGINT,
    nm_id BIGINT,
    brand TEXT,
    status INTEGER,
    order_uid TEXT REFERENCES orders(order_uid) ON DELETE CASCADE
);

-- Индексы
CREATE INDEX idx_orders_customer_id ON orders(customer_id);
CREATE INDEX idx_items_order_uid ON items(order_uid);
CREATE INDEX idx_orders_track_number ON orders(track_number);
