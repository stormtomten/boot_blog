package main

import (
	"boot_dev/boot_blog/internal/config"
	"boot_dev/boot_blog/internal/database"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type command struct {
	name string
	args []string
}

type commands struct {
	handlers map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	f, ok := c.handlers[cmd.name]
	if !ok {
		return fmt.Errorf("unknown command: %v", cmd.name)
	}
	return f(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) error {
	if c.handlers == nil {
		c.handlers = make(map[string]func(*state, command) error)
	}
	if _, exists := c.handlers[name]; exists {
		return fmt.Errorf("command %q already registered", name)
	}
	c.handlers[name] = f
	return nil
}

type state struct {
	cfg *config.Config
	db  *database.Queries
}

func main() {
	cfg, err := config.ReadConfig()
	if err != nil {
		log.Fatalf("Config error: %v", err)
	}

	db, err := sql.Open("postgres", cfg.DbUrl)
	if err != nil {
		log.Fatalf("Database error: %v\n", err)
	}

	dbQueries := database.New(db)
	ps := state{cfg: &cfg, db: dbQueries}

	c := commands{handlers: make(map[string]func(*state, command) error)}

	if err := c.register("login", handlerLogin); err != nil {
		log.Fatalf("error: %v\n", err)
	}

	if err := c.register("register", handlerRegister); err != nil {
		log.Fatalf("error: %v\n", err)
	}
	if err := c.register("reset", handlerReset); err != nil {
		log.Fatalf("error: %v\n", err)
	}
	if err := c.register("users", handlerGetUsers); err != nil {
		log.Fatalf("error: %v\n", err)
	}
	if err := c.register("agg", handlerAggregate); err != nil {
		log.Fatalf("error: %v\n", err)
	}
	if err := c.register("addfeed", middlewareLoggedIn(handlerCreateFeed)); err != nil {
		log.Fatalf("error: %v\n", err)
	}
	if err := c.register("feeds", handlerGetFeeds); err != nil {
		log.Fatalf("error: %v\n", err)
	}
	if err := c.register("follow", middlewareLoggedIn(handlerCreateFeedFollow)); err != nil {
		log.Fatalf("error: %v\n", err)
	}
	if err := c.register("following", middlewareLoggedIn(handlerGetFeedFollowsForUser)); err != nil {
		log.Fatalf("error: %v\n", err)
	}
	if err := c.register("unfollow", middlewareLoggedIn(handlerUnFollow)); err != nil {
		log.Fatalf("error: %v\n", err)
	}
	if err := c.register("browse", middlewareLoggedIn(handlerBrowse)); err != nil {
		log.Fatalf("error: %v\n", err)
	}
	args := os.Args
	if len(args) < 2 {
		log.Fatal("Usage: cli <command> [args...]")
		os.Exit(1)
	}
	cmd := command{name: args[1], args: args[2:]}
	if err := c.run(&ps, cmd); err != nil {
		log.Fatalf("fatal error: %v\n", err)
		os.Exit(1)
	}

	os.Exit(0)
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: gator login <username>")
	}
	uname := cmd.args[0]
	ctx := context.Background()
	if _, err := s.db.GetUser(ctx, uname); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("user %s does not exist", uname)
		}
		return fmt.Errorf("database error: %v", err)
	}
	if err := s.cfg.SetUser(uname); err != nil {
		return fmt.Errorf("failed to set user in config: %v", err)
	}
	fmt.Printf("User set to %s\n", uname)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: gator register <username>")
	}
	uname := cmd.args[0]
	ctx := context.Background()
	_, err := s.db.GetUser(ctx, uname)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("database error: %v", err)
	}
	if err == nil {
		return fmt.Errorf("user %s already exists", uname)
	}
	_, err = s.db.CreateUser(ctx, database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      uname,
	})
	if err != nil {
		return fmt.Errorf("database error: %v", err)
	}
	if err := s.cfg.SetUser(uname); err != nil {
		return fmt.Errorf("failed to set user in config: %v", err)
	}
	fmt.Printf("User %s created successfully\n", uname)

	return nil
}

func handlerGetUsers(s *state, cmd command) error {
	if len(cmd.args) > 0 {
		return fmt.Errorf("usage: gator users")
	}
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("database error: %v", err)
	}
	if len(users) < 1 {
		fmt.Printf("no users: %v\n", users)
		return nil
	}
	for _, u := range users {
		fmt.Printf("* %s%s\n", u, func() string {
			if u == s.cfg.CurrentUserName {
				return " (current)"
			}
			return ""
		}())
	}

	return nil
}

func handlerReset(s *state, cmd command) error {
	if len(cmd.args) > 0 {
		return fmt.Errorf("usage: gator reset")
	}
	if err := s.db.DeleteAllUsers(context.Background()); err != nil {
		return fmt.Errorf("database error: %v", err)
	}
	fmt.Println("Database reset successful.")
	return nil
}

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			return err
		}
		err = handler(s, cmd, user)
		return err
	}
}
