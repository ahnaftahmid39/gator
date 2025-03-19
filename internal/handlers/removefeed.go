package handlers

import (
	"context"
	"fmt"

	"github.com/ahnaftahmid39/gator/internal/app"
	"github.com/ahnaftahmid39/gator/internal/database"
)

func HandleRemoveFeed(s *app.State, cmd app.Command, user database.User) error {
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
