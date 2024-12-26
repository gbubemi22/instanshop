-- 000001_create_users_table.up.sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    verified_email BOOLEAN DEFAULT FALSE,
    otp_token VARCHAR(255),
    expired_at TIMESTAMP,
    role VARCHAR(50) NOT NULL DEFAULT 'user',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    CONSTRAINT idx_deleted_at_deleted_at_idx UNIQUE(deleted_at)
);
