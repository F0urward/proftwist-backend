package repository

const (
	queryGetAllRoadmap = `
        SELECT id, owner_id, category_id, name, description, is_public, color, 
               referenced_roadmap_id, subscriber_count, created_at, updated_at 
        FROM roadmaps`
	queryGetByID = `
        SELECT id, owner_id, category_id, name, description, is_public, color, 
               referenced_roadmap_id, subscriber_count, created_at, updated_at 
        FROM roadmaps 
        WHERE id = $1`
	queryCreate = `
        INSERT INTO roadmaps 
        (id, owner_id, category_id, name, description, is_public, color, referenced_roadmap_id, subscriber_count, created_at, updated_at) 
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	queryUpdate = `
        UPDATE roadmaps 
        SET category_id = $2, name = $3, description = $4, is_public = $5, 
            color = $6, referenced_roadmap_id = $7, updated_at = $8 
        WHERE id = $1`
	queryDelete = `
        DELETE FROM roadmaps 
        WHERE id = $1`
)
