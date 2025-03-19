package main

import (
	"database/sql"

	// "encoding/json"
	"fmt"
	"os"

	"github.com/ahnaftahmid39/gator/internal/app"
	"github.com/ahnaftahmid39/gator/internal/config"
	"github.com/ahnaftahmid39/gator/internal/database"
	"github.com/ahnaftahmid39/gator/internal/handlers"
	"github.com/ahnaftahmid39/gator/internal/middlewares"
)

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

	// init app.State with config and database query
	dbQueries := database.New(db)
	s := &app.State{
		Cfg: cfg,
		Db:  dbQueries,
	}

	// create commands
	cmds := app.Commands{
		CmdHandlerMap: make(map[string]func(*app.State, app.Command) error),
	}
	cmds.Register("addfeed", middlewares.LoggedIn(handlers.HandleAddFeed))
	cmds.Register("agg", handlers.HandleAggregator)
	cmds.Register("browse", middlewares.LoggedIn(handlers.HandleBrowse))
	cmds.Register("feeds", handlers.HandleFeeds)
	cmds.Register("follow", middlewares.LoggedIn(handlers.HandleFollow))
	cmds.Register("following", middlewares.LoggedIn(handlers.HandleFollowing))
	cmds.Register("login", handlers.HandleLogin)
	cmds.Register("register", handlers.HandleRegister)
	cmds.Register("reset", handlers.HandleReset)
	cmds.Register("users", handlers.HandleUsers)
	cmds.Register("unfollow", middlewares.LoggedIn(handlers.HandleUnfollow))
	cmds.Register("removefeed", middlewares.LoggedIn(handlers.HandleRemoveFeed))

	// handle command
	if len(os.Args) < 2 {
		fmt.Printf("No command name given\n")
		os.Exit(1)
	}

	cmd := app.Command{
		Name: os.Args[1],
		Args: os.Args[2:],
	}
	err = cmds.Run(s, cmd)

	if err != nil {
		fmt.Printf("Error executing %s: %v\n", cmd.Name, err)
		os.Exit(1)
	}

}
