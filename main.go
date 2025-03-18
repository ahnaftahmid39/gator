package main

import (
	"context"
	"database/sql"
	"encoding/json"

	// "encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/ahnaftahmid39/gator/internal/config"
	"github.com/ahnaftahmid39/gator/internal/database"
	"github.com/ahnaftahmid39/gator/internal/rss"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type state struct {
	cfg *config.Config
	db  *database.Queries
}

type command struct {
	name string
	args []string
}

type commands struct {
	cmdHandlerMap map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.cmdHandlerMap[name] = f
}
func (c *commands) run(s *state, cmd command) error {
	if handler, exists := c.cmdHandlerMap[cmd.name]; exists {
		err := handler(s, cmd)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("command does not exist")
	}
	return nil
}

func handleReset(s *state, cmd command) error {
	ctx := context.Background()
	err := s.db.DeleteAllUsers(ctx)
	return err
}

func handleUsers(s *state, cmd command) error {
	ctx := context.Background()
	users, err := s.db.GetUsers(ctx)
	if err != nil {
		return err
	}
	for _, user := range users {
		if user.Name == s.cfg.CurrentUserName {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}

	return nil
}

func handleRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("the register handler expects a single argument, the username")
	}
	userName := cmd.args[0]

	ctx := context.Background()
	user, err := s.db.CreateUser(ctx, database.CreateUserParams{
		ID:   uuid.New(),
		Name: userName,
		CreatedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		UpdatedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	})

	if err != nil {
		return err
	}

	err = s.cfg.SetUser(userName)
	if err != nil {
		return err
	}

	pretty, _ := json.MarshalIndent(user, "", "  ")
	fmt.Printf("the user has been created and set in config. User:\n%+v\n", string(pretty))
	return nil
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("the login handler expects a single argument, the username")
	}

	userName := cmd.args[0]
	ctx := context.Background()
	user, err := s.db.GetUserByName(ctx, userName)
	if err != nil {
		return fmt.Errorf("user does not exist in the databse. please use register")
	}
	err = s.cfg.SetUser(userName)
	if err != nil {
		return err
	}

	fmt.Printf("the user %s has been set\n", user.Name)
	return nil
}

func handleAggregator(s *state, cmd command) error {
	ctx := context.Background()
	feed, err := rss.FetchFeed(ctx, "https://www.wagslane.dev/index.xml")

	if err != nil {
		return err
	}

	fmt.Println(feed)

	return nil
}

func handleAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 2 {
		return fmt.Errorf("not enough arguments, command syntax: addfeed <feed_name> <feed_url>")
	}

	ctx := context.Background()
	feed, err := s.db.CreateFeed(ctx, database.CreateFeedParams{
		ID:   uuid.New(),
		Name: cmd.args[0],
		Url:  cmd.args[1],
		CreatedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		UpdatedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		UserID: user.ID,
	})

	if err != nil {
		return fmt.Errorf("error creating a new feed, %w", err)
	}

	feed_follow, err := s.db.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
		CreatedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		UpdatedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		UserID: user.ID,
		FeedID: feed.ID,
	})

	if err != nil {
		return fmt.Errorf("error creating a feed follow, %w", err)
	}

	pretty, _ := json.MarshalIndent(feed, "", "  ")
	fmt.Printf("The Feed has been created. Feed:\n%+v\n", string(pretty))
	fmt.Printf("A Feed follow has been created. FeedName:%s, User: %s\n", feed_follow.FeedName, feed_follow.UserName)

	return nil

}

func handleFeeds(s *state, cmd command) error {
	ctx := context.Background()
	feeds, err := s.db.GetAllFeeds(ctx)
	if err != nil {
		return fmt.Errorf("error getting all the feeds, %w", err)
	}
	for _, feed := range feeds {
		fmt.Println("Feed Name:", feed.Name, "Feed URL:", feed.Url, "Feed Created By:", feed.UserName)
	}

	return nil
}

func handleFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("follow command expects a feed url")
	}

	ctx := context.Background()
	feed, err := s.db.GetFeedByUrl(ctx, cmd.args[0])
	if err != nil {
		return fmt.Errorf("error while getting feed by url, %w", err)
	}

	feed_follow, err := s.db.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
		CreatedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		UpdatedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		UserID: user.ID,
		FeedID: feed.ID,
	})

	if err != nil {
		return fmt.Errorf("error while creating feed follow, %w", err)
	}

	fmt.Println("Feed Name:", feed_follow.FeedName, ",User:", feed_follow.UserName)

	return nil
}

func handleFollowing(s *state, cmd command, user database.User) error {
	ctx := context.Background()
	feed_follows, err := s.db.GetFeedFollowsForUser(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("error while getting feed follows, %w", err)
	}

	for _, follow := range feed_follows {
		fmt.Println("*", follow.FeedName)
	}

	return nil
}

func middlewareLoggedIn(
	handler func(s *state, cmd command, user database.User) error,
) func(*state, command) error {
	return func(s *state, cmd command) error {
		ctx := context.Background()
		user, err := s.db.GetUserByName(ctx, s.cfg.CurrentUserName)
		if err != nil {
			return fmt.Errorf("error while getting current user information, %w", err)
		}

		return handler(s, cmd, user)
	}

}

func main() {
	// read config
	cfg, err := config.Read()
	if err != nil {
		fmt.Println("Error loading config:", err)
		os.Exit(1)
	}

	// connect to database
	db, err := sql.Open("postgres", cfg.DBUrl)
	if err != nil {
		fmt.Println("Error opening database:", err)
		os.Exit(1)
	}

	// init state with config and database query
	dbQueries := database.New(db)
	s := &state{
		cfg: cfg,
		db:  dbQueries,
	}

	// create commands
	cmds := commands{
		cmdHandlerMap: make(map[string]func(*state, command) error),
	}
	cmds.register("login", handlerLogin)
	cmds.register("register", handleRegister)
	cmds.register("reset", handleReset)
	cmds.register("users", handleUsers)
	cmds.register("agg", handleAggregator)
	cmds.register("addfeed", middlewareLoggedIn(handleAddFeed))
	cmds.register("feeds", handleFeeds)
	cmds.register("follow", middlewareLoggedIn(handleFollow))
	cmds.register("following", middlewareLoggedIn(handleFollowing))

	// handle command
	if len(os.Args) < 2 {
		fmt.Printf("No command name given\n")
		os.Exit(1)
	}

	cmd := command{
		name: os.Args[1],
		args: os.Args[2:],
	}
	err = cmds.run(s, cmd)

	if err != nil {
		fmt.Printf("Error executing %s: %v\n", cmd.name, err)
		os.Exit(1)
	}

}
