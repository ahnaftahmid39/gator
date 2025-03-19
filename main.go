package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"strconv"

	// "encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/ahnaftahmid39/gator/internal/app"
	"github.com/ahnaftahmid39/gator/internal/config"
	"github.com/ahnaftahmid39/gator/internal/database"
	"github.com/ahnaftahmid39/gator/internal/rss"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type commands struct {
	cmdHandlerMap map[string]func(*app.State, app.Command) error
}

func (c *commands) register(name string, f func(*app.State, app.Command) error) {
	c.cmdHandlerMap[name] = f
}
func (c *commands) run(s *app.State, cmd app.Command) error {
	if handler, exists := c.cmdHandlerMap[cmd.Name]; exists {
		err := handler(s, cmd)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("command does not exist")
	}
	return nil
}

func handleReset(s *app.State, cmd app.Command) error {
	ctx := context.Background()
	err := s.Db.DeleteAllUsers(ctx)
	return err
}

func handleUsers(s *app.State, cmd app.Command) error {
	ctx := context.Background()
	users, err := s.Db.GetUsers(ctx)
	if err != nil {
		return err
	}
	for _, user := range users {
		if user.Name == s.Cfg.CurrentUserName {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}

	return nil
}

func handleRegister(s *app.State, cmd app.Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("the register handler expects a single argument, the username")
	}
	userName := cmd.Args[0]

	ctx := context.Background()
	user, err := s.Db.CreateUser(ctx, database.CreateUserParams{
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

	err = s.Cfg.SetUser(userName)
	if err != nil {
		return err
	}

	pretty, _ := json.MarshalIndent(user, "", "  ")
	fmt.Printf("the user has been created and set in config. User:\n%+v\n", string(pretty))
	return nil
}

func handlerLogin(s *app.State, cmd app.Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("the login handler expects a single argument, the username")
	}

	userName := cmd.Args[0]
	ctx := context.Background()
	user, err := s.Db.GetUserByName(ctx, userName)
	if err != nil {
		return fmt.Errorf("user does not exist in the databse. please use register")
	}
	err = s.Cfg.SetUser(userName)
	if err != nil {
		return err
	}

	fmt.Printf("the user %s has been set\n", user.Name)
	return nil
}

func scrapeFeeds(s *app.State) error {
	ctx := context.Background()
	feedToFetch, err := s.Db.GetNextFeedToFetch(ctx)
	if err != nil {
		return fmt.Errorf("error retrieving next feed to fetch, %w", err)
	}

	fmt.Printf("----Fetchin %s-----\n", feedToFetch.Name)
	feed, err := rss.FetchFeed(ctx, feedToFetch.Url)
	if err != nil {
		return err
	}

	_, err = s.Db.MarkFeedFetchedById(ctx, database.MarkFeedFetchedByIdParams{
		ID: feedToFetch.ID,
		LastFetchedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	})

	if err != nil {
		return fmt.Errorf("error marking feed fetched by id, %w", err)
	}

	for _, item := range feed.Channel.Item {
		fmt.Println("*", item.Title)
		publishDate, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			return fmt.Errorf("error parsing publish date, %w", err)
		}
		_, err = s.Db.CreatePost(ctx, database.CreatePostParams{
			ID: uuid.New(),
			CreatedAt: sql.NullTime{
				Time:  time.Now(),
				Valid: true,
			},
			UpdatedAt: sql.NullTime{
				Time:  time.Now(),
				Valid: true,
			},
			Title: sql.NullString{
				String: item.Title,
				Valid:  true,
			},
			Description: sql.NullString{
				String: item.Description,
				Valid:  true,
			},
			Url: item.Link,
			PublishedAt: sql.NullTime{
				Time:  publishDate,
				Valid: true,
			},
			FeedID: feedToFetch.ID,
		})

		if err != nil {
			if pqErr, ok := err.(*pq.Error); !ok || pqErr.Code != "23505" {
				return fmt.Errorf("error creating post, %w", err)
			}
		}
	}
	fmt.Println("-------------------")
	return nil
}

func handleAggregator(s *app.State, cmd app.Command) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("missing arguments. syntax: agg <time_between_reqs>")
	}

	duration, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("error parsing time_between_reqs. Valid examples: 1s, 10h10m1s, 10m etc, Error: %w", err)
	}

	t := time.NewTicker(duration)
	for ; ; <-t.C {
		err = scrapeFeeds(s)
		if err != nil {
			return err
		}
	}
}

func handleAddFeed(s *app.State, cmd app.Command, user database.User) error {
	if len(cmd.Args) < 2 {
		return fmt.Errorf("not enough arguments, app.Command syntax: addfeed <feed_name> <feed_url>")
	}

	ctx := context.Background()
	feed, err := s.Db.CreateFeed(ctx, database.CreateFeedParams{
		ID:   uuid.New(),
		Name: cmd.Args[0],
		Url:  cmd.Args[1],
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

	feed_follow, err := s.Db.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
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

func handleRemoveFeed(s *app.State, cmd app.Command, user database.User) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("missing arguments. Syntax: removefeed <feed_url>")
	}

	deletedFeed, err := s.Db.DeleteFeedByUrlAndUser(context.Background(), database.DeleteFeedByUrlAndUserParams{
		Url:    cmd.Args[0],
		UserID: user.ID,
	})

	if err != nil {
		return fmt.Errorf("error deleting feed by url and user_id, you are likely trying to remove a feed that isn't yours, %w", err)
	}

	fmt.Println("Successfully deleted feed", deletedFeed.Name)
	return nil
}

func handleFeeds(s *app.State, cmd app.Command) error {
	ctx := context.Background()
	feeds, err := s.Db.GetAllFeeds(ctx)
	if err != nil {
		return fmt.Errorf("error getting all the feeds, %w", err)
	}
	for _, feed := range feeds {
		fmt.Println("Feed Name:", feed.Name, "Feed URL:", feed.Url, "Feed Created By:", feed.UserName)
	}

	return nil
}

func handleFollow(s *app.State, cmd app.Command, user database.User) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("follow command expects a feed url")
	}

	ctx := context.Background()
	feed, err := s.Db.GetFeedByUrl(ctx, cmd.Args[0])
	if err != nil {
		return fmt.Errorf("error while getting feed by url, %w", err)
	}

	feed_follow, err := s.Db.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
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

func handleFollowing(s *app.State, cmd app.Command, user database.User) error {
	ctx := context.Background()
	feed_follows, err := s.Db.GetFeedFollowsForUser(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("error while getting feed follows, %w", err)
	}

	for _, follow := range feed_follows {
		fmt.Println("*", follow.FeedName)
	}

	return nil
}

func handleUnfollow(s *app.State, cmd app.Command, user database.User) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("not enough arguments. Need feed_url. Syntax: unfollow <feed_url>")
	}

	err := s.Db.DeleteFeedFollowByFeedUrlAndUserId(context.Background(), database.DeleteFeedFollowByFeedUrlAndUserIdParams{
		Url:    cmd.Args[0],
		UserID: user.ID,
	})

	if err != nil {
		return fmt.Errorf("error deleting feed follow by feed url and userid, %w", err)
	}

	fmt.Println("Successfully unfollowed feed")

	return nil
}

func handleBrowse(s *app.State, cmd app.Command, user database.User) error {

	limit := 2
	if len(cmd.Args) > 0 {
		l, err := strconv.Atoi(cmd.Args[0])
		if err != nil {
			return fmt.Errorf("error parsing argument limit, %w", err)
		}
		limit = l
	}
	posts, err := s.Db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})

	if err != nil {
		return fmt.Errorf("error getting posts for user, %w", err)
	}

	for i, post := range posts {
		fmt.Printf("%v. ", i+1)
		if post.Title.Valid {
			fmt.Println(post.Title.String)
		}
		fmt.Println(post.Url)
		if post.Description.Valid {
			fmt.Println(post.Description.String)
		}
		fmt.Println()
	}

	return nil

}

func middlewareLoggedIn(
	handler func(s *app.State, cmd app.Command, user database.User) error,
) func(*app.State, app.Command) error {
	return func(s *app.State, cmd app.Command) error {
		ctx := context.Background()
		user, err := s.Db.GetUserByName(ctx, s.Cfg.CurrentUserName)
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

	// init app.State with config and database query
	dbQueries := database.New(db)
	s := &app.State{
		Cfg: cfg,
		Db:  dbQueries,
	}

	// create commands
	cmds := commands{
		cmdHandlerMap: make(map[string]func(*app.State, app.Command) error),
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
	cmds.register("unfollow", middlewareLoggedIn(handleUnfollow))
	cmds.register("removefeed", middlewareLoggedIn(handleRemoveFeed))
	cmds.register("browse", middlewareLoggedIn(handleBrowse))

	// handle command
	if len(os.Args) < 2 {
		fmt.Printf("No command name given\n")
		os.Exit(1)
	}

	cmd := app.Command{
		Name: os.Args[1],
		Args: os.Args[2:],
	}
	err = cmds.run(s, cmd)

	if err != nil {
		fmt.Printf("Error executing %s: %v\n", cmd.Name, err)
		os.Exit(1)
	}

}
