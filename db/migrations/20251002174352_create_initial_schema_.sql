-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TYPE user_role AS ENUM ('admin', 'regular');
CREATE TYPE chat_type AS ENUM ('direct', 'group');
CREATE TYPE member_role AS ENUM ('owner', 'admin', 'member');

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
                              roadmap_id TEXT NOT NULL,
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

CREATE TABLE group_chat (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title TEXT,
    avatar_url TEXT,
    roadmap_node_id TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE direct_chat (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user1_id UUID NOT NULL REFERENCES "user"(id),
    user2_id UUID NOT NULL REFERENCES "user"(id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(user1_id, user2_id)
);

CREATE TABLE group_chat_members (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    group_chat_id UUID NOT NULL REFERENCES group_chat(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES "user"(id),
    UNIQUE(group_chat_id, user_id)
);

CREATE TABLE group_chat_messages (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    group_chat_id UUID NOT NULL REFERENCES group_chat(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES "user"(id),
    content TEXT NOT NULL,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE direct_chat_messages (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    direct_chat_id UUID NOT NULL REFERENCES direct_chat(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES "user"(id),
    content TEXT NOT NULL,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS direct_chat_messages;
DROP TABLE IF EXISTS group_chat_messages;
DROP TABLE IF EXISTS group_chat_members;
DROP TABLE IF EXISTS direct_chat;
DROP TABLE IF EXISTS group_chat;
DROP TABLE IF EXISTS chats;
DROP TABLE IF EXISTS roadmap_info_subscription;
DROP TABLE IF EXISTS roadmap_info;
DROP TABLE IF EXISTS category;
DROP TABLE IF EXISTS vk_user;
DROP TABLE IF EXISTS "user";

DROP TYPE IF EXISTS message_type;
DROP TYPE IF EXISTS member_role;
DROP TYPE IF EXISTS chat_type;
DROP TYPE IF EXISTS user_role;
-- +goose StatementEnd
