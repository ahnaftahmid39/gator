package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/ahnaftahmid39/gator/internal/app"
	"github.com/ahnaftahmid39/gator/internal/database"
	"github.com/ahnaftahmid39/gator/internal/rss"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

func scrapeFeeds(s *app.State) error {
	ctx := context.Background()
	feedToFetch, err := s.Db.GetNextFeedToFetch(ctx)
	if err != nil {
		return fmt.Errorf("error retrieving next feed to fetch, %w", err)
	}

	fmt.Printf("----Fetchin %s-----\n", feedToFetch.Name)
	feed, err := rss.FetchFeed(ctx, feedToFetch.Url)
	if err != nil {
		return err
	}

	_, err = s.Db.MarkFeedFetchedById(ctx, database.MarkFeedFetchedByIdParams{
		ID: feedToFetch.ID,
		LastFetchedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	})

	if err != nil {
		return fmt.Errorf("error marking feed fetched by id, %w", err)
	}

	for _, item := range feed.Channel.Item {
		fmt.Println("*", item.Title)
		publishDate, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			return fmt.Errorf("error parsing publish date, %w", err)
		}
		_, err = s.Db.CreatePost(ctx, database.CreatePostParams{
			ID: uuid.New(),
			CreatedAt: sql.NullTime{
				Time:  time.Now(),
				Valid: true,
			},
			UpdatedAt: sql.NullTime{
				Time:  time.Now(),
				Valid: true,
			},
			Title: sql.NullString{
				String: item.Title,
				Valid:  true,
			},
			Description: sql.NullString{
				String: item.Description,
				Valid:  true,
			},
			Url: item.Link,
			PublishedAt: sql.NullTime{
				Time:  publishDate,
				Valid: true,
			},
			FeedID: feedToFetch.ID,
		})

		if err != nil {
			if pqErr, ok := err.(*pq.Error); !ok || pqErr.Code != "23505" {
				return fmt.Errorf("error creating post, %w", err)
			}
		}
	}
	fmt.Println("-------------------")
	return nil
}

func HandleAggregator(s *app.State, cmd app.Command) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("missing arguments. syntax: agg <time_between_reqs>")
	}

	duration, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("error parsing time_between_reqs. Valid examples: 1s, 10h10m1s, 10m etc, Error: %w", err)
	}

	t := time.NewTicker(duration)
	for ; ; <-t.C {
		err = scrapeFeeds(s)
		if err != nil {
			return err
		}
	}
}
