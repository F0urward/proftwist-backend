package repository

const (
	queryCreateChat = `
		INSERT INTO chats (id, type, title, description, avatar_url, created_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	queryGetChat = `
		SELECT id, type, title, description, avatar_url, created_by, created_at, updated_at
		FROM chats 
		WHERE id = $1`

	queryGetUserChats = `
		SELECT c.id, c.type, c.title, c.description, c.avatar_url, c.created_by, c.created_at, c.updated_at
		FROM chats c
		INNER JOIN chat_members cm ON c.id = cm.chat_id
		WHERE cm.user_id = $1
		ORDER BY c.updated_at DESC`

	queryUpdateChatTimestamp = `
		UPDATE chats SET updated_at = $1 WHERE id = $2`

	querySaveMessage = `
		INSERT INTO messages (id, chat_id, user_id, content, metadata, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`

	queryGetChatMessages = `
		SELECT id, chat_id, user_id, content, metadata, created_at, updated_at
		FROM messages 
		WHERE chat_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	queryAddChatMember = `
		INSERT INTO chat_members (chat_id, user_id, role, joined_at, last_read)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (chat_id, user_id) DO UPDATE SET
			role = EXCLUDED.role`

	queryRemoveChatMember = `
		DELETE FROM chat_members WHERE chat_id = $1 AND user_id = $2`

	queryIsChatMember = `
		SELECT 1 FROM chat_members WHERE chat_id = $1 AND user_id = $2`

	queryGetChatMembers = `
		SELECT id, chat_id, user_id, role, joined_at, last_read
		FROM chat_members
		WHERE chat_id = $1
		ORDER BY joined_at`

	queryDeleteChat = `
        DELETE FROM chats WHERE id = $1`

	queryDeleteChatMesseges = `DELETE FROM messages WHERE chat_id = $1`

	queryDeleteChatMembers = `DELETE FROM chat_members WHERE chat_id = $1`
)
