package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/ahnaftahmid39/gator/internal/app"
	"github.com/ahnaftahmid39/gator/internal/database"
)

func HandleFollow(s *app.State, cmd app.Command, user database.User) error {
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
