package repository

const (
	queryGetAll = `
		SELECT id, name, description, created_at, updated_at 
		FROM category
		ORDER BY created_at DESC`

	queryGetByID = `
		SELECT id, name, description, created_at, updated_at 
		FROM category 
		WHERE id = $1`

	queryGetByName = `
		SELECT id, name, description, created_at, updated_at 
		FROM category 
		WHERE name = $1`

	queryCreate = `
		INSERT INTO category 
		(name, description) 
		VALUES ($1, $2)
		RETURNING id, name, description, created_at, updated_at `

	queryUpdate = `
		UPDATE category 
		SET name = $2, description = $3, updated_at = $4 
		WHERE id = $1`

	queryDelete = `
		DELETE FROM category 
		WHERE id = $1`
)
