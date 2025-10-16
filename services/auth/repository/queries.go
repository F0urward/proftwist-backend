package repository

const (
	queryCreateUser = `
		INSERT INTO "user" 
		(username, email, password_hash, role, avatar_url) 
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at`

	queryGetUserByEmail = `
		SELECT id, username, email, password_hash, role, avatar_url, created_at, updated_at 
		FROM "user" 
		WHERE email = $1`

	queryGetUserByID = `
		SELECT id, username, email, password_hash, role, avatar_url, created_at, updated_at 
		FROM "user" 
		WHERE id = $1`

	queryUpdateUser = `
		UPDATE "user" 
		SET username = $2, email = $3, password_hash = $4, role = $5, avatar_url = $6, updated_at = $7 
		WHERE id = $1`

	queryDeleteUser = `
		DELETE FROM "user" 
		WHERE id = $1`

	queryCreateVKUser = `
		INSERT INTO vk_user 
		(user_id, vk_user_id, access_token, refresh_token, expires_at, device_id) 
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at`

	queryGetVKUserByUserID = `
		SELECT id, user_id, vk_user_id, access_token, refresh_token, expires_at, device_id, created_at, updated_at 
		FROM vk_user 
		WHERE user_id = $1`

	queryGetVKUserByVKUserID = `
		SELECT id, user_id, vk_user_id, access_token, refresh_token, expires_at, device_id, created_at, updated_at 
		FROM vk_user 
		WHERE vk_user_id = $1`

	queryUpdateVKUser = `
		UPDATE vk_user 
		SET vk_user_id = $2, access_token = $3, refresh_token = $4, expires_at = $5, device_id = $6, updated_at = NOW() 
		WHERE id = $1`

	queryDeleteVKUser = `
		DELETE FROM vk_user 
		WHERE user_id = $1`
)
