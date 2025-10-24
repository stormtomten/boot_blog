-- +goose Up
CREATE TABLE posts (
    id UUID UNIQUE PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    title TEXT ,
    url TEXT UNIQUE NOT NULL,
    description TEXT ,
    published_at TIMESTAMP,
    feed_id UUID NOT NULL REFERENCES feeds(id)
    );
CREATE INDEX idx_posts_feed_published
ON posts (feed_id, published_at DESC);

-- +goose Down
DROP TABLE posts;
