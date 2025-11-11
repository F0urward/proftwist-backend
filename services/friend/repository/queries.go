package repository

const (
	queryCreateFriendship = `
		INSERT INTO friends (user_id, friend_id)
		VALUES ($1, $2)`

	queryDeleteFriendship = `
		DELETE FROM friends 
		WHERE user_id = $1 AND friend_id = $2`

	queryGetFriendIDs = `
		SELECT friend_id 
		FROM friends 
		WHERE user_id = $1`

	queryIsFriends = `
		SELECT 1 FROM friends 
		WHERE user_id = $1 AND friend_id = $2`

	queryCreateFriendRequest = `
		INSERT INTO friend_requests (from_user_id, to_user_id, message)
		VALUES ($1, $2, $3)
		RETURNING id, status, created_at, updated_at`

	queryGetFriendRequestByID = `
		SELECT id, from_user_id, to_user_id, status, message, created_at, updated_at
		FROM friend_requests 
		WHERE id = $1`

	queryGetFriendRequestsForUser = `
		SELECT id, from_user_id, to_user_id, status, message, created_at, updated_at
		FROM friend_requests 
		WHERE to_user_id = $1 AND status = 'pending'`

	queryGetSentFriendRequests = `
		SELECT id, from_user_id, to_user_id, status, message, created_at, updated_at
		FROM friend_requests 
		WHERE from_user_id = $1 AND status = 'pending'`

	queryUpdateFriendRequestStatus = `
		UPDATE friend_requests 
		SET status = $1, updated_at = NOW()
		WHERE id = $2`

	queryDeleteFriendRequest = `
		DELETE FROM friend_requests 
		WHERE id = $1`

	queryGetFriendRequestBetweenUsers = `
		SELECT id, from_user_id, to_user_id, status, message, created_at, updated_at
		FROM friend_requests 
		WHERE (from_user_id = $1 AND to_user_id = $2) OR (from_user_id = $2 AND to_user_id = $1)`

	queryGetPendingFriendRequestBetweenUsers = `
		SELECT id, from_user_id, to_user_id, status, message, created_at, updated_at
		FROM friend_requests 
		WHERE ((from_user_id = $1 AND to_user_id = $2) OR (from_user_id = $2 AND to_user_id = $1))
		AND status = 'pending'`
)
