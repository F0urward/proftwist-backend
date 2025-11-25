package repository

const (
	queryCreateMaterial = `
		INSERT INTO materials (name, url, roadmap_node_id, author_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at`

	queryGetMaterialByID = `
		SELECT id, name, url, roadmap_node_id, author_id, created_at, updated_at
		FROM materials 
		WHERE id = $1`

	queryGetMaterialsByNode = `
		SELECT id, name, url, roadmap_node_id, author_id, created_at, updated_at
		FROM materials 
		WHERE roadmap_node_id = $1
		ORDER BY created_at DESC`

	queryGetMaterialsByAuthor = `
		SELECT id, name, url, roadmap_node_id, author_id, created_at, updated_at
		FROM materials 
		WHERE author_id = $1
		ORDER BY created_at DESC`

	queryDeleteMaterial = `
		DELETE FROM materials WHERE id = $1`
)
