-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TYPE user_role AS ENUM ('admin', 'regular');

CREATE TABLE "user" (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT,
    role user_role DEFAULT 'regular',
    avatar_url TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE vk_user (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL UNIQUE REFERENCES "user"(id) ON DELETE CASCADE,
    vk_user_id TEXT NOT NULL UNIQUE,
    access_token TEXT NOT NULL,
    refresh_token TEXT,
    expires_at TIMESTAMP NOT NULL,
    device_id TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE category (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE roadmap_info (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    roadmap_id TEXT NOT NULL DEFAULT encode(gen_random_bytes(12), 'hex'),
    author_id UUID NOT NULL REFERENCES "user"(id) ON DELETE CASCADE,
    category_id UUID REFERENCES category(id),
    name TEXT NOT NULL,
    description TEXT,
    is_public BOOLEAN DEFAULT true,
    referenced_roadmap_info_id UUID REFERENCES roadmap_info(id),
    subscriber_count INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE roadmap_info_subscription (
    user_id UUID REFERENCES "user"(id) ON DELETE CASCADE,
    roadmap_info_id UUID REFERENCES roadmap_info(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (user_id, roadmap_info_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS roadmap_info_subscription;
DROP TABLE IF EXISTS roadmap_info;
DROP TABLE IF EXISTS category;
DROP TABLE IF EXISTS vk_user;
DROP TABLE IF EXISTS "user";

DROP TYPE IF EXISTS user_role;
-- +goose StatementEnd
