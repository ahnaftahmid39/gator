package handlers

import (
	"context"
	"fmt"

	"github.com/ahnaftahmid39/gator/internal/app"
)

func HandleFeeds(s *app.State, cmd app.Command) error {
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
