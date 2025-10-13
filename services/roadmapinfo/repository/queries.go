package repository

const (
	queryGetAll = `
		SELECT id, roadmap_id, author_id, name, description, is_public, 
               referenced_roadmap_info_id, subscriber_count, created_at, updated_at 
        FROM roadmap_info`

	queryGetByID = `
        SELECT id, roadmap_id, author_id, name, description, is_public,
               referenced_roadmap_info_id, subscriber_count, created_at, updated_at 
        FROM roadmap_info 
        WHERE id = $1`

	queryGetByRoadmapID = `
        SELECT id, roadmap_id, author_id, name, description, is_public,
               referenced_roadmap_info_id, subscriber_count, created_at, updated_at 
        FROM roadmap_info 
        WHERE roadmap_id = $1`

	queryCreate = `
        INSERT INTO roadmap_info 
        (author_id, name, description, is_public, referenced_roadmap_info_id, roadmap_id, subscriber_count) 
        VALUES ($1, $2, $3, $4, $5, $6, $7)`

	queryUpdate = `
        UPDATE roadmap_info 
        SET name = $2, description = $3, is_public = $4, 
            referenced_roadmap_info_id = $5, roadmap_id = $6, updated_at = $7 
        WHERE id = $1`

	queryDelete = `
        DELETE FROM roadmap_info 
        WHERE id = $1`
)
