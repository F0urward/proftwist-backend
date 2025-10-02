# proftwist-backend


CREATE TYPE user_role AS ENUM ('admin', 'regular');

CREATE TYPE node_type AS ENUM ('root', 'topic', 'leaf');

CREATE TABLE app_user (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(100) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role user_role DEFAULT 'regular',
    avatar_url TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE category (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    color VARCHAR(7),
    icon VARCHAR(50),
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE roadmap (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    owner_id UUID NOT NULL REFERENCES "user"(id) ON DELETE CASCADE,
    category_id UUID NOT NULL REFERENCES category(id),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    is_public BOOLEAN DEFAULT true,
    color VARCHAR(7),
    referenced_roadmap_id UUID REFERENCES roadmap(id),
    subscriber_count INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE node (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    roadmap_id UUID NOT NULL REFERENCES roadmap(id) ON DELETE CASCADE,
    title VARCHAR(500) NOT NULL,
    description TEXT,
    type node_type NOT NULL,
    x_position INTEGER,
    y_position INTEGER,
    width INTEGER,
    height INTEGER,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE roadmap_subscription (
    user_id UUID REFERENCES "user"(id) ON DELETE CASCADE,
    roadmap_id UUID REFERENCES roadmap(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (user_id, roadmap_id)
);

DROP TABLE IF EXISTS roadmap_subscription;
DROP TABLE IF EXISTS node;
DROP TABLE IF EXISTS roadmap;
DROP TABLE IF EXISTS category;
DROP TABLE IF EXISTS app_user;

DROP TYPE IF EXISTS node_type;
DROP TYPE IF EXISTS user_role;
