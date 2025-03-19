package handlers

import (
	"context"

	"github.com/ahnaftahmid39/gator/internal/app"
)

func HandleReset(s *app.State, cmd app.Command) error {
	ctx := context.Background()
	err := s.Db.DeleteAllUsers(ctx)
	return err
}
