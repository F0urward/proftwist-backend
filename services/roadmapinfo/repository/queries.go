package repository

const (
	queryGetAll = `
		SELECT id, author_id, category_id, name, description, is_public, 
               referenced_roadmap_info_id, subscriber_count, created_at, updated_at 
        FROM roadmap_info`

	queryGetByID = `
        SELECT id, author_id, category_id, name, description, is_public,
               referenced_roadmap_info_id, subscriber_count, created_at, updated_at 
        FROM roadmap_info 
        WHERE id = $1`

	queryCreate = `
        INSERT INTO roadmap_info 
        (author_id, category_id, name, description, is_public, referenced_roadmap_info_id, subscriber_count) 
        VALUES ($1, $2, $3, $4, $5, $6, $7)`

	queryUpdate = `
        UPDATE roadmap_info 
        SET category_id = $2, name = $3, description = $4, is_public = $5, referenced_roadmap_info_id = $6, updated_at = $7 
        WHERE id = $1`

	queryDelete = `
        DELETE FROM roadmap_info 
        WHERE id = $1`
)
