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
)
