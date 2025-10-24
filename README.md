# Gator - RSS Feed Aggregator

A command-line RSS feed aggregator built in Go, following the [Boot.dev course](https://www.boot.dev/courses/build-blog-aggregator-golang).

Gator allows you to aggregate RSS feeds, manage users, follow feeds, and browse posts through a simple CLI interface.

## Features

- **User Management**: Register and login users
- **Feed Aggregation**: Automatically fetch and aggregate posts from RSS feeds
- **Feed Following**: Follow/unfollow RSS feeds
- **Post Browsing**: Browse recent posts from followed feeds
- **Database Storage**: PostgreSQL-backed with type-safe queries using sqlc

## Installation

### Prerequisites

- Go 1.19 or later
- PostgreSQL database

### Build from Source

```bash
git clone <repository-url>
cd boot_blog
go build ./cmd/gator
```

This will create a `gator` executable in the current directory.

## Database Setup

1. Create a PostgreSQL database
2. Run the database migrations (using goose):

```bash
# Install goose if you don't have it
go install github.com/pressly/goose/v3/cmd/goose@latest

# Run migrations
goose -dir sql/schema postgres "your-connection-string" up
```

## Configuration

Create a configuration file at `~/.gatorconfig.json`:

```json
{
  "db_url": "postgres://username:password@localhost:5432/gator_db?sslmode=disable"
}
```

## Usage

### Getting Started

1. Register a new user:
```bash
./gator register <username>
```

2. Login with your user:
```bash
./gator login <username>
```

### Available Commands

#### User Management
- `register <username>` - Create a new user account
- `login <username>` - Login as an existing user
- `users` - List all registered users

#### Feed Management
- `addfeed <name> <url>` - Add a new RSS feed (requires login)
- `feeds` - List all available feeds
- `follow <url>` - Follow an existing feed (requires login)
- `following` - List feeds you're following (requires login)
- `unfollow <url>` - Unfollow a feed (requires login)

#### Content Aggregation
- `agg <duration>` - Start the feed aggregator (e.g., `agg 30s` for 30-second intervals)
- `browse [limit]` - Browse recent posts from followed feeds (default limit: 2, requires login)

#### Utilities
- `reset` - Delete all users and reset the database

### Examples

```bash
# Register and login
./gator register alice
./gator login alice

# Add and follow feeds
./gator addfeed "Boot.dev Blog" "https://blog.boot.dev/rss"
./gator follow "https://blog.boot.dev/rss"

# Start aggregating feeds every 5 minutes
./gator agg 5m

# Browse recent posts (in another terminal)
./gator browse 5
```

## Dependencies

- [github.com/google/uuid](https://github.com/google/uuid) - UUID generation
- [github.com/lib/pq](https://github.com/lib/pq) - PostgreSQL driver
- [sqlc](https://sqlc.dev/) - Type-safe SQL code generation
- [goose](https://github.com/pressly/goose) - Database migrations

## Development

### Code Generation

After modifying SQL queries in `sql/queries/`, regenerate the database code:

```bash
sqlc generate
```

### Testing

```bash
go test ./...
```

### Linting

```bash
go vet ./...
gofmt -l .
```

## Limitations

- No post read/unread status tracking
- No post deletion capabilities
- No web interface (CLI only)
- Basic error handling for RSS parsing edge cases

## Contributing

This project was built following the Boot.dev Go course.

## License

This project is for educational purposes as part of the Boot.dev curriculum.
