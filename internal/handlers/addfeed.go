package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ahnaftahmid39/gator/internal/app"
	"github.com/ahnaftahmid39/gator/internal/database"
	"github.com/google/uuid"
)

func HandleAddFeed(s *app.State, cmd app.Command, user database.User) error {
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
