package models

import (
	"time"

	"github.com/MazzMS/go-rss/internal/database"
	"github.com/google/uuid"
)

type FeedFollow struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    uuid.UUID `json:"user_id"`
	FeedID    uuid.UUID `json:"feed_id"`
}

func DBFeedFollowToFeedFollow(dbFF database.FeedsFollow) FeedFollow {
	return FeedFollow{
		ID:        dbFF.ID,
		CreatedAt: dbFF.CreatedAt,
		UpdatedAt: dbFF.UpdatedAt,
		UserID:    dbFF.UserID,
		FeedID:    dbFF.FeedID,
	}
}

func DBFeedsFollowsToFeedsFollows(dbFFs []database.FeedsFollow) []FeedFollow {
	feedsFollows := make([]FeedFollow, len(dbFFs))
	for i, dbFF := range dbFFs {
		feedsFollows[i] = DBFeedFollowToFeedFollow(dbFF)
	}
	return feedsFollows
}
