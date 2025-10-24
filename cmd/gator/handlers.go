package main

import (
	"boot_dev/boot_blog/internal/database"
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
)

func handlerCreateFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 2 {
		return fmt.Errorf("usage: gator addfeed <name> <url>")
	}
	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
		Url:       cmd.args[1],
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("database error: %v", err)
	}
	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		Url:       cmd.args[1],
	})
	fmt.Println("Feed created successfully:")
	printFeed(feed)
	fmt.Println()
	fmt.Println("=====================================")

	return nil
}

func printFeed(feed database.Feed) {
	// fmt.Printf("* ID:            %s\n", feed.ID)
	fmt.Printf("* Created:       %v\n", feed.CreatedAt)
	// fmt.Printf("* Updated:       %v\n", feed.UpdatedAt)
	fmt.Printf("* Name:          %s\n", feed.Name)
	fmt.Printf("* URL:           %s\n", feed.Url)
	// fmt.Printf("* UserID:        %s\n", feed.UserID)
}

func handlerGetFeeds(s *state, cmd command) error {
	if len(cmd.args) > 0 {
		return fmt.Errorf("usage: gator feeds")
	}
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("database error: %v", err)
	}
	for _, feed := range feeds {
		fmt.Printf("Name: %s\n", feed.Name)
		fmt.Printf("URL: %s\n", feed.Url)
		fmt.Printf("User: %s\n", feed.Username.String)
		fmt.Println("====")
	}

	return nil
}

func handlerCreateFeedFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: gator follow <url>")
	}
	feed, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		Url:       cmd.args[0],
	})
	if err != nil {
		return fmt.Errorf("database error: %v", err)
	}
	feed_name, err := s.db.GetFeedName(context.Background(), feed.FeedID)
	if err != nil {
		return fmt.Errorf("database error: %v", err)
	}
	fmt.Println("Follow created successfully:")
	fmt.Printf("Name: %v\nUser: %s\n", feed_name, feed.UserName)
	fmt.Println("=====================================")

	return nil
}

func handlerGetFeedFollowsForUser(s *state, cmd command, user database.User) error {
	if len(cmd.args) > 0 {
		return fmt.Errorf("usage: gator following")
	}
	feeds, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("database error: %v", err)
	}
	for _, feed := range feeds {
		fmt.Printf("Name: %s\n", feed.FeedName)
		fmt.Printf("User: %s\n", feed.UserName)
		fmt.Println("====")
	}

	return nil
}

func handlerUnFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: gator follow <url>")
	}
	err := s.db.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{
		Url:    cmd.args[0],
		UserID: user.ID,
	})
	if err != nil {
		return fmt.Errorf("database error: %v", err)
	}
	fmt.Println("Unfollowed successfully:")

	return nil
}

func handlerBrowse(s *state, cmd command, user database.User) error {
	if len(cmd.args) > 1 {
		return fmt.Errorf("usage: gator browse <limit>")
	}
	limit := int32(2)
	if len(cmd.args) >= 1 {
		if n64, err := strconv.ParseInt(cmd.args[0], 10, 32); err == nil {
			limit = int32(n64)
		}
	}
	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  limit,
	})
	if err != nil {
		return fmt.Errorf("database error: %v", err)
	}
	for _, post := range posts {
		fmt.Printf("Feed: %s | Title: %s| Description: %s\n", post.Name, post.Title.String, post.Description.String)
	}
	return nil
}
