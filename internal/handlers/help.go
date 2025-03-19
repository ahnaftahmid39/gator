package handlers

import (
	"fmt"

	"github.com/ahnaftahmid39/gator/internal/app"
)

func HandleHelp(s *app.State, cmd app.Command) error {
	fmt.Println("Available Commands:")
	fmt.Println("  gator register <username>")
	fmt.Println("  gator login <username>")
	fmt.Println("  gator users")
	fmt.Println("  gator addfeed <feed_name> <feed_url>")
	fmt.Println("  gator feeds")
	fmt.Println("  gator follow <feed_url>")
	fmt.Println("  gator following")
	fmt.Println("  gator unfollow <feed_url>")
	fmt.Println("  gator agg <time_between_requests>")
	fmt.Println("  gator browse [limit]")

	return nil
}
