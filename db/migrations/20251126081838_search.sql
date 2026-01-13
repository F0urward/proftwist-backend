-- +goose Up
-- +goose StatementBegin

CREATE TEXT SEARCH CONFIGURATION russian_fts (COPY = simple);
ALTER TEXT SEARCH CONFIGURATION russian_fts
    ALTER MAPPING FOR asciiword, asciihword, hword_asciipart, word, hword, hword_part
    WITH simple;

-- Roadmap info FTS setup
ALTER TABLE roadmap_info 
ADD COLUMN fts tsvector;

UPDATE roadmap_info 
SET fts = 
    setweight(to_tsvector('russian_fts', coalesce(name, '')), 'A') ||
    setweight(to_tsvector('russian_fts', coalesce(description, '')), 'B');

CREATE INDEX roadmap_info_fts_idx ON roadmap_info USING GIN(fts);

CREATE INDEX roadmap_info_fts_public_idx ON roadmap_info USING GIN(fts) 
WHERE is_public = true;

-- User FTS setup
ALTER TABLE "user" 
ADD COLUMN fts tsvector;

UPDATE "user" 
SET fts = setweight(to_tsvector('russian_fts', coalesce(username, '')), 'A');

CREATE INDEX user_fts_idx ON "user" USING GIN(fts);

-- Common functions
CREATE OR REPLACE FUNCTION make_prefix_tsquery(query_text text)
RETURNS tsquery AS $$
BEGIN
    IF query_text IS NULL OR trim(query_text) = '' THEN
        RETURN to_tsquery('');
    END IF;
    
    RETURN (
        SELECT string_agg(lexeme || ':*', ' & ')::tsquery
        FROM unnest(to_tsvector('simple', query_text))
    );
END;
$$ LANGUAGE plpgsql IMMUTABLE;

-- Roadmap info trigger function
CREATE OR REPLACE FUNCTION roadmap_info_fts_update()
RETURNS trigger AS $$
BEGIN
    NEW.fts = 
        setweight(to_tsvector('russian_fts', coalesce(NEW.name, '')), 'A') ||
        setweight(to_tsvector('russian_fts', coalesce(NEW.description, '')), 'B');
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- User trigger function
CREATE OR REPLACE FUNCTION user_fts_update()
RETURNS trigger AS $$
BEGIN
    NEW.fts = setweight(to_tsvector('russian_fts', coalesce(NEW.username, '')), 'A');
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create triggers
CREATE TRIGGER roadmap_info_fts_trigger
    BEFORE INSERT OR UPDATE OF name, description ON roadmap_info
    FOR EACH ROW EXECUTE FUNCTION roadmap_info_fts_update();

CREATE TRIGGER user_fts_trigger
    BEFORE INSERT OR UPDATE OF username ON "user"
    FOR EACH ROW EXECUTE FUNCTION user_fts_update();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Drop triggers
DROP TRIGGER IF EXISTS roadmap_info_fts_trigger ON roadmap_info;
DROP TRIGGER IF EXISTS user_fts_trigger ON "user";

-- Drop functions
DROP FUNCTION IF EXISTS roadmap_info_fts_update();
DROP FUNCTION IF EXISTS user_fts_update();
DROP FUNCTION IF EXISTS make_prefix_tsquery;

-- Drop indexes
DROP INDEX IF EXISTS roadmap_info_fts_public_idx;
DROP INDEX IF EXISTS roadmap_info_fts_idx;
DROP INDEX IF EXISTS user_fts_idx;

-- Drop FTS columns
ALTER TABLE roadmap_info DROP COLUMN IF EXISTS fts;
ALTER TABLE "user" DROP COLUMN IF EXISTS fts;

DROP TEXT SEARCH CONFIGURATION IF EXISTS russian_fts;

-- +goose StatementEnd
