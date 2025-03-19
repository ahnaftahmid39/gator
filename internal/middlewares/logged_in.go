package middlewares

import (
	"context"
	"fmt"

	"github.com/ahnaftahmid39/gator/internal/app"
	"github.com/ahnaftahmid39/gator/internal/database"
)

func LoggedIn(
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
