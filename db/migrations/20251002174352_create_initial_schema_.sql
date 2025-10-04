-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE user_role AS ENUM ('admin', 'regular');

CREATE TABLE "user" (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    role user_role DEFAULT 'regular',
    avatar_url TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE category (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL,
    description TEXT,
    icon TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE roadmap_info (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    author_id UUID NOT NULL REFERENCES "user"(id) ON DELETE CASCADE,
    category_id UUID NOT NULL REFERENCES category(id),
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
DROP TABLE IF EXISTS "user";

DROP TYPE IF EXISTS user_role;
-- +goose StatementEnd
