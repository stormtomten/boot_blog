-- name: CreatePost :one
INSERT INTO posts(id, created_at, updated_at, title, url, description, published_at , feed_id)
VALUES(
    $1,
    NOW(),
    NOW(),
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: GetPostsForUser :many
SELECT posts.*, feeds.name
FROM feed_follows
INNER JOIN feeds ON feeds.id = feed_follows.feed_id
INNER JOIN posts ON feeds.id = posts.feed_id
WHERE feed_follows.user_id = $1
ORDER BY posts.published_at DESC, posts.id DESC
LIMIT $2;

