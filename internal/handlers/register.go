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

func HandleRegister(s *app.State, cmd app.Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("the register handler expects a single argument, the username")
	}
	userName := cmd.Args[0]

	ctx := context.Background()
	user, err := s.Db.CreateUser(ctx, database.CreateUserParams{
		ID:   uuid.New(),
		Name: userName,
		CreatedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		UpdatedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	})

	if err != nil {
		return err
	}

	err = s.Cfg.SetUser(userName)
	if err != nil {
		return err
	}

	pretty, _ := json.MarshalIndent(user, "", "  ")
	fmt.Printf("the user has been created and set in config. User:\n%+v\n", string(pretty))
	return nil
}
