package handlers

import (
	"context"
	"fmt"

	"github.com/ahnaftahmid39/gator/internal/app"
)

func HandleLogin(s *app.State, cmd app.Command) error {
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
