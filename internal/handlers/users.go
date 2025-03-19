package handlers

import (
	"context"
	"fmt"

	"github.com/ahnaftahmid39/gator/internal/app"
)

func HandleUsers(s *app.State, cmd app.Command) error {
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
