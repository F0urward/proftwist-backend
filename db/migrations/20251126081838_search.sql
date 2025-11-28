-- +goose Up
-- +goose StatementBegin

CREATE TEXT SEARCH CONFIGURATION russian_fts (COPY = simple);
ALTER TEXT SEARCH CONFIGURATION russian_fts
    ALTER MAPPING FOR asciiword, asciihword, hword_asciipart, word, hword, hword_part
    WITH simple;

ALTER TABLE roadmap_info 
ADD COLUMN fts tsvector;

UPDATE roadmap_info 
SET fts = 
    setweight(to_tsvector('russian_fts', coalesce(name, '')), 'A') ||
    setweight(to_tsvector('russian_fts', coalesce(description, '')), 'B');

CREATE INDEX roadmap_info_fts_idx ON roadmap_info USING GIN(fts);

CREATE INDEX roadmap_info_fts_public_idx ON roadmap_info USING GIN(fts) 
WHERE is_public = true;

CREATE OR REPLACE FUNCTION roadmap_info_fts_update()
RETURNS trigger AS $$
BEGIN
    NEW.fts = 
        setweight(to_tsvector('russian_fts', coalesce(NEW.name, '')), 'A') ||
        setweight(to_tsvector('russian_fts', coalesce(NEW.description, '')), 'B');
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER roadmap_info_fts_trigger
    BEFORE INSERT OR UPDATE OF name, description ON roadmap_info
    FOR EACH ROW EXECUTE FUNCTION roadmap_info_fts_update();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TRIGGER IF EXISTS roadmap_info_fts_trigger ON roadmap_info;
DROP FUNCTION IF EXISTS roadmap_info_fts_update();

DROP INDEX IF EXISTS roadmap_info_fts_public_idx;
DROP INDEX IF EXISTS roadmap_info_fts_idx;

ALTER TABLE roadmap_info DROP COLUMN IF EXISTS fts;

DROP TEXT SEARCH CONFIGURATION IF EXISTS russian_fts;

-- +goose StatementEnd
