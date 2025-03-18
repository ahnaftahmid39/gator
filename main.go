package main

import (
	"context"
	"database/sql"

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
	_, err := s.db.CreateUser(ctx, database.CreateUserParams{
		ID: uuid.NullUUID{
			UUID:  uuid.New(),
			Valid: true,
		},
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

	// pretty, _ := json.MarshalIndent(user, "", "  ")
	// fmt.Printf("the user has been created and set in config. User:\n%+v\n", string(pretty))
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
