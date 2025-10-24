-- name: CreateFeed :one
INSERT INTO feeds(id, created_at, updated_at, name, url, user_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: GetFeeds :many
SELECT
feeds.name AS name,
feeds.url,
users.name AS username
FROM feeds
LEFT JOIN users ON feeds.user_id = users.id;

-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS(
    INSERT INTO feed_follows( id, created_at, updated_at, user_id, feed_id)
    VALUES(
        $1,
        $2,
        $3,
        $4,
        (SELECT id FROM feeds WHERE feeds.url = $5)
    )RETURNING *
)
SELECT
    inserted_feed_follow.*,
    feeds.name AS feed_name,
    users.name AS user_name
FROM inserted_feed_follow
INNER JOIN feeds ON inserted_feed_follow.feed_id = feeds.id
INNER JOIN users ON inserted_feed_follow.user_id = users.id;

-- name: GetFeedName :one
SELECT name
FROM feeds
WHERE id = $1;

-- name: GetFeedFollowsForUser :many
SELECT
    feeds.name AS feed_name,
    users.name AS user_name
FROM feed_follows
INNER JOIN feeds ON feed_follows.feed_id = feeds.id
INNER JOIN users ON feed_follows.user_id = users.id
WHERE users.id = $1;

-- name: DeleteFeedFollow :exec
DELETE
FROM feed_follows
WHERE feed_id = (SELECT id 
    FROM feeds
    WHERE url = $1)
AND feed_follows.user_id = $2;

-- name: MarkFeedFetched :exec
UPDATE feeds
SET last_fetched_at = NOW(), updated_at = NOW()
WHERE id = $1;

-- name: GetNextFeedToFetch :one
SELECT id, url, name
FROM feeds
ORDER BY last_fetched_at ASC NULLS FIRST, id ASC
LIMIT 1;
