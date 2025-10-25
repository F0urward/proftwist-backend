package repository

const (
	queryGetAll = `
        SELECT id, roadmap_id, author_id, category_id, name, description, is_public, 
               referenced_roadmap_info_id, subscriber_count, created_at, updated_at 
        FROM roadmap_info`

	queryGetAllByCategoryID = `
        SELECT id, roadmap_id, author_id, category_id, name, description, is_public,
               referenced_roadmap_info_id, subscriber_count, created_at, updated_at 
        FROM roadmap_info 
        WHERE category_id = $1`

	queryGetByID = `
        SELECT id, roadmap_id, author_id, category_id, name, description, is_public,
               referenced_roadmap_info_id, subscriber_count, created_at, updated_at 
        FROM roadmap_info 
        WHERE id = $1`

	queryGetByRoadmapID = `
        SELECT id, roadmap_id, author_id, category_id, name, description, is_public,
               referenced_roadmap_info_id, subscriber_count, created_at, updated_at 
        FROM roadmap_info 
        WHERE roadmap_id = $1`

	queryCreate = `
        INSERT INTO roadmap_info 
        (author_id, category_id, name, description, is_public, referenced_roadmap_info_id, roadmap_id, subscriber_count) 
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
        RETURNING id, roadmap_id, author_id, category_id, name, description, is_public,
               referenced_roadmap_info_id, subscriber_count, created_at, updated_at`

	queryUpdate = `
        UPDATE roadmap_info 
        SET category_id = $2, name = $3, description = $4, is_public = $5, 
            referenced_roadmap_info_id = $6, roadmap_id = $7, updated_at = $8 
        WHERE id = $1`

	queryDelete = `
        DELETE FROM roadmap_info 
        WHERE id = $1`
)
