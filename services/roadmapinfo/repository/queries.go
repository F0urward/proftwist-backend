package repository

const (
	queryGetAllPublic = `
        SELECT id, roadmap_id, author_id, category_id, name, description, is_public,
               referenced_roadmap_info_id, created_at, updated_at 
        FROM roadmap_info 
        WHERE is_public = true`

	queryGetAllPublicByCategoryID = `
        SELECT id, roadmap_id, author_id, category_id, name, description, is_public,
                referenced_roadmap_info_id, created_at, updated_at 
        FROM roadmap_info 
        WHERE category_id = $1 AND is_public = true`

	queryGetAllByUserID = `
        SELECT id, roadmap_id, author_id, category_id, name, description, is_public,
        referenced_roadmap_info_id, created_at, updated_at 
        FROM roadmap_info 
        WHERE author_id = $1`

	queryGetByID = `
        SELECT id, roadmap_id, author_id, category_id, name, description, is_public,
               referenced_roadmap_info_id, created_at, updated_at 
        FROM roadmap_info 
        WHERE id = $1`

	queryGetByRoadmapID = `
        SELECT id, roadmap_id, author_id, category_id, name, description, is_public,
               referenced_roadmap_info_id, created_at, updated_at 
        FROM roadmap_info 
        WHERE roadmap_id = $1`

	queryCreate = `
        INSERT INTO roadmap_info 
        (author_id, category_id, name, description, is_public, referenced_roadmap_info_id, roadmap_id) 
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING id, roadmap_id, author_id, category_id, name, description, is_public,
               referenced_roadmap_info_id, created_at, updated_at`

	queryUpdate = `
        UPDATE roadmap_info 
        SET category_id = $2, name = $3, description = $4, is_public = $5, 
            referenced_roadmap_info_id = $6, roadmap_id = $7, updated_at = $8 
        WHERE id = $1`

	queryDelete = `
        DELETE FROM roadmap_info 
        WHERE id = $1`

	queryCreateSubscription = `
        INSERT INTO roadmap_info_subscription (user_id, roadmap_info_id)
        VALUES ($1, $2)
        ON CONFLICT (user_id, roadmap_info_id) DO NOTHING`

	queryDeleteSubscription = `
        DELETE FROM roadmap_info_subscription 
        WHERE user_id = $1 AND roadmap_info_id = $2`

	querySubscriptionExists = `
        SELECT EXISTS(
            SELECT 1 FROM roadmap_info_subscription 
            WHERE user_id = $1 AND roadmap_info_id = $2
        )`

	queryGetSubscribedRoadmapIDs = `
        SELECT roadmap_info_id 
        FROM roadmap_info_subscription 
        WHERE user_id = $1`

	queryGetByIDs = `
        SELECT id, roadmap_id, author_id, category_id, name, description, is_public,
               referenced_roadmap_info_id, created_at, updated_at 
        FROM roadmap_info 
        WHERE id = ANY($1)`

	querySearchPublic = `
        SELECT id, roadmap_id, author_id, category_id, name, description, is_public,
               referenced_roadmap_info_id, created_at, updated_at 
        FROM roadmap_info 
        WHERE fts @@ plainto_tsquery('simple', $1)
        AND is_public = true`

	querySearchPublicWithCategory = `
        SELECT id, roadmap_id, author_id, category_id, name, description, is_public,
               referenced_roadmap_info_id, created_at, updated_at 
        FROM roadmap_info 
        WHERE fts @@ plainto_tsquery('simple', $1)
        AND is_public = true
        AND category_id = $2`
)
