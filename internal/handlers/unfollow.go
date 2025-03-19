package handlers

import (
	"context"
	"fmt"

	"github.com/ahnaftahmid39/gator/internal/app"
	"github.com/ahnaftahmid39/gator/internal/database"
)

func HandleUnfollow(s *app.State, cmd app.Command, user database.User) error {
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
