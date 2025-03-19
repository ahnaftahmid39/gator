package handlers

import (
	"context"
	"fmt"

	"github.com/ahnaftahmid39/gator/internal/app"
	"github.com/ahnaftahmid39/gator/internal/database"
)

func HandleFollowing(s *app.State, cmd app.Command, user database.User) error {
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
