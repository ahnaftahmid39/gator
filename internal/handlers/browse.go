package handlers

import (
	"context"
	"fmt"
	"strconv"

	"github.com/ahnaftahmid39/gator/internal/app"
	"github.com/ahnaftahmid39/gator/internal/database"
)

func HandleBrowse(s *app.State, cmd app.Command, user database.User) error {

	limit := 2
	if len(cmd.Args) > 0 {
		l, err := strconv.Atoi(cmd.Args[0])
		if err != nil {
			return fmt.Errorf("error parsing argument limit, %w", err)
		}
		limit = l
	}
	posts, err := s.Db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})

	if err != nil {
		return fmt.Errorf("error getting posts for user, %w", err)
	}

	for i, post := range posts {
		fmt.Printf("%v. ", i+1)
		if post.Title.Valid {
			fmt.Println(post.Title.String)
		}
		fmt.Println(post.Url)
		if post.Description.Valid {
			fmt.Println(post.Description.String)
		}
		fmt.Println()
	}

	return nil

}
