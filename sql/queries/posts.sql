-- name: CreatePost :one
INSERT INTO posts (id, created_at, updated_at, title, url, description, publised_at, feed_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetPostsByUser :many
SELECT
posts.id,
posts.created_at,
posts.updated_at,
posts.title,
posts.url,
posts.description,
posts.publised_at
FROM posts
JOIN feeds
    ON posts.feed_id = feeds.id
JOIN users
    ON users.id = feeds.user_id
WHERE users.id = $1
ORDER BY posts.publised_at
LIMIT $2;