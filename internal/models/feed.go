package models

import (
	"time"

	"github.com/MazzMS/go-rss/internal/database"
	"github.com/google/uuid"
)

type Feed struct {
	ID            uuid.UUID `json:"id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Name          string    `json:"name"`
	Url           string    `json:"url"`
	UserID        uuid.UUID `json:"user_id"`
	LastFetchedAt time.Time `json:"last_fetched_at"`
}

func DBFeedToFeed(dbFeed database.Feed) Feed {
	return Feed{
		ID:            dbFeed.ID,
		CreatedAt:     dbFeed.CreatedAt,
		UpdatedAt:     dbFeed.UpdatedAt,
		Name:          dbFeed.Name,
		Url:           dbFeed.Url,
		UserID:        dbFeed.UserID,
		LastFetchedAt: dbFeed.LastFetchedAt.Time,
	}
}

func DBFeedsToFeeds(dbFeeds []database.Feed) []Feed {
	feeds := make([]Feed, len(dbFeeds))
	for i, dbFeed := range dbFeeds {
		feeds[i] = DBFeedToFeed(dbFeed)
	}
	return feeds
}
