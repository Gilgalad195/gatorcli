package main

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/Gilgalad195/gatorcli/internal/database"
	"github.com/Gilgalad195/gatorcli/internal/webconn"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

func scrapeFeeds(s *state) error {
	ctx := context.Background()

	nextFeed, err := s.db.GetNextFeedToFetch(ctx)
	if err != nil {
		return err
	}

	toBeMarked := database.MarkFeedFetchedParams{
		ID: nextFeed.ID,
		LastFetchedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		UpdatedAt: time.Now(),
	}

	if err := s.db.MarkFeedFetched(ctx, toBeMarked); err != nil {
		return err
	}

	feed, err := webconn.FetchFeed(ctx, nextFeed.Url)
	if err != nil {
		return err
	}

	for _, item := range feed.Channel.Item {
		desc := sql.NullString{}
		if item.Description != "" {
			desc = sql.NullString{String: item.Description, Valid: true}
		}
		pub := sql.NullTime{}
		if t, ok := parseMaybeTime(item.PubDate); ok {
			pub = sql.NullTime{Time: t, Valid: true}
		}

		postParams := database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       item.Title,
			Url:         item.Link,
			Description: desc,
			PublishedAt: pub,
			FeedID:      nextFeed.ID,
		}

		_, err := s.db.CreatePost(ctx, postParams)
		if err != nil {
			if pgErr, ok := err.(*pq.Error); ok {
				if pgErr.Code == "23505" {
					continue
				}
			}
			log.Printf("an error occured: %v", err)
			continue
		}
	}
	return nil
}
