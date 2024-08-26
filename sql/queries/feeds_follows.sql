-- name: CreateFeedFollow :one
INSERT INTO feeds_follows (id, created_at, updated_at, user_id, feed_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: DeleteFeedFollow :exec
DELETE FROM feeds_follows WHERE user_id = $1 and id = $2;

-- name: GetFeedsFollowsByUser :many
SELECT * FROM feeds_follows WHERE user_id = $1;
