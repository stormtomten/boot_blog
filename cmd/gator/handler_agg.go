package main

import (
	"boot_dev/boot_blog/internal/database"
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
)

func handlerAggregate(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("incorrect arguments, usage: gator agg <timeduration>")
	}
	timeBetweenRequests, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return fmt.Errorf("failed to parse time: %v", err)
	}
	log.Printf("Collecting feeds every %s...", timeBetweenRequests)
	ticker := time.NewTicker(timeBetweenRequests)
	defer ticker.Stop()
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}

	return nil
}

func scrapeFeeds(s *state) error {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("database error: %v", err)
	}

	err = s.db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		return fmt.Errorf("database error: %v", err)
	}

	feedData, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		return fmt.Errorf("request error: %v", err)
	}
	for _, post := range feedData.Channel.Item {
		_, err := s.db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			Title:       parseString(post.Title),
			Url:         post.Link,
			Description: parseString(post.Description),
			PublishedAt: parseRSSTime(post.PubDate),
			FeedID:      feed.ID,
		})
		if err != nil {
			log.Printf("failed to create post: %v\n", err)
		}
	}

	log.Printf("Feed %s collected, %v posts found", feed.Name, len(feedData.Channel.Item))

	return nil
}

func parseRSSTime(s string) sql.NullTime {
	s = strings.TrimSpace(s)
	if s == "" {
		return sql.NullTime{Valid: false}
	}
	layouts := []string{
		time.RFC1123Z, time.RFC1123, time.RFC822Z, time.RFC822,
		time.RFC3339, time.RFC3339Nano, time.RubyDate,
	}
	for _, l := range layouts {
		if t, err := time.Parse(l, s); err == nil {
			return sql.NullTime{Time: t.UTC(), Valid: true}
		}
	}
	return sql.NullTime{Valid: false}
}

func parseString(s string) sql.NullString {
	s = strings.TrimSpace(s)
	if s == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: s, Valid: true}
}
