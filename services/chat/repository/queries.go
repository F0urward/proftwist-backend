package repository

const (
	queryGetGroupChatByNode = `
		SELECT id, title, avatar_url, roadmap_node_id, created_at, updated_at
		FROM group_chat 
		WHERE roadmap_node_id = $1`

	queryGetGroupChatsByUser = `
		SELECT gc.id, gc.title, gc.avatar_url, gc.roadmap_node_id, gc.created_at, gc.updated_at
		FROM group_chat gc
		INNER JOIN group_chat_members gcm ON gc.id = gcm.group_chat_id
		WHERE gcm.user_id = $1
		ORDER BY gc.updated_at DESC`

	queryGetGroupChatMembers = `
		SELECT id, group_chat_id, user_id
		FROM group_chat_members
		WHERE group_chat_id = $1`

	queryIsGroupChatMember = `
		SELECT 1 FROM group_chat_members WHERE group_chat_id = $1 AND user_id = $2`

	queryAddGroupChatMember = `
		INSERT INTO group_chat_members (group_chat_id, user_id)
		VALUES ($1, $2)`

	queryRemoveGroupChatMember = `
		DELETE FROM group_chat_members WHERE group_chat_id = $1 AND user_id = $2`

	queryGeDirectChatsByUser = `
		SELECT id, user1_id, user2_id, created_at, updated_at
		FROM direct_chat 
		WHERE user1_id = $1 OR user2_id = $2
		ORDER BY updated_at DESC`

	queryGetDirectChat = `
		SELECT id, user1_id, user2_id, created_at, updated_at
		FROM direct_chat 
		WHERE id = $1`

	queryIsDirectChatMember = `
		SELECT 1 FROM direct_chat WHERE id = $1 AND (user1_id = $2 OR user2_id = $2)`

	querySaveGroupMessage = `
		INSERT INTO group_chat_messages (id, group_chat_id, user_id, content, metadata, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`

	queryGetGroupChatMessages = `
		SELECT id, group_chat_id, user_id, content, metadata, created_at, updated_at
		FROM group_chat_messages 
		WHERE group_chat_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	querySaveDirectMessage = `
		INSERT INTO direct_chat_messages (id, direct_chat_id, user_id, content, metadata, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`

	queryGetDirectChatMessages = `
		SELECT id, direct_chat_id, user_id, content, metadata, created_at, updated_at
		FROM direct_chat_messages 
		WHERE direct_chat_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`
)
