-- +goose Up
-- +goose StatementBegin
ALTER TABLE group_chat_messages ADD COLUMN thread_root_id UUID REFERENCES group_chat_messages(id) ON DELETE CASCADE;
CREATE INDEX idx_group_chat_messages_thread_root_id ON group_chat_messages(thread_root_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_group_chat_messages_thread_root_id;
ALTER TABLE group_chat_messages DROP COLUMN IF EXISTS thread_root_id;
-- +goose StatementEnd
