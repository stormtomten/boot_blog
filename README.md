# boot_blog/gator

Following this [course](https://www.boot.dev/courses/build-blog-aggregator-golang "Build a Blog Aggregator in Go") by the folks over at [boot.dev](https://www.boot.dev/), I present you with **gator** a RSS aggregator.

## What it does

It provides functions for fetching RSS feeds, aggregates posts into a database and provides CLI commands for user management, feed following, and browsing posts.

## What it doesn't do

Won't let you read, mark as read and/or delete posts.

```json
{
  "db_url": "protocol://username:password@host:port/database?sslmode=disable"
}
```

.gatorconfig.json
