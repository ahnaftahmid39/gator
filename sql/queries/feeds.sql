-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: GetAllFeeds :many
SELECT f.name name, f.url url, u.name user_name
FROM feeds f INNER JOIN users u ON u.id = f.user_id;

-- name: GetFeedByUrl :one
SELECT * from feeds
WHERE url = $1;


-- name: GetNextFeedToFetch :one
SELECT id, name, url
FROM feeds
ORDER BY last_fetched_at ASC NULLS FIRST
LIMIT 1;

-- name: MarkFeedFetchedById :one
UPDATE feeds
SET
    last_fetched_at=$2,
    updated_at=$2
WHERE id = $1
RETURNING *;


-- name: DeleteFeedByUrlAndUser :one
DELETE
FROM feeds
WHERE url = $1 AND user_id = $2
RETURNING *;
