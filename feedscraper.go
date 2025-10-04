package main

import (
	"context"
	"fmt"

	"github.com/Gilgalad195/gatorcli/internal/database"
	"github.com/Gilgalad195/gatorcli/internal/webconn"
)

func scrapeFeeds(s *state) error {
	ctx := context.Background()

	nextFeed, err := s.db.GetNextFeedToFetch(ctx)
	if err != nil {
		return err
	}

	toBeMarked := database.MarkFeedFetchedParams{
		ID:            nextFeed.ID,
		LastFetchedAt: nextFeed.LastFetchedAt,
		UpdatedAt:     nextFeed.UpdatedAt,
	}

	if err := s.db.MarkFeedFetched(ctx, toBeMarked); err != nil {
		return err
	}

	feed, err := webconn.FetchFeed(ctx, nextFeed.Url)
	if err != nil {
		return err
	}

	for _, item := range feed.Channel.Item {
		fmt.Printf("%v\n", item.Title)
	}
	return nil
}
