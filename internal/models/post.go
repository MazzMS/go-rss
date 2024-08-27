package models

import (
	"time"

	"github.com/MazzMS/go-rss/internal/database"
	"github.com/google/uuid"
)

type Post struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Title       string    `json:"title"`
	Url         string    `json:"url"`
	Description string    `json:"description"`
	PublishedAt time.Time `json:"publised_at"`
	FeedID      uuid.UUID `json:"feed_id"`
}

func DBPostToPost(dbPost database.Post) Post {
	return Post{
		ID:          dbPost.ID,
		CreatedAt:   dbPost.CreatedAt,
		UpdatedAt:   dbPost.UpdatedAt,
		PublishedAt: dbPost.PublishedAt,
		Title:       dbPost.Title,
		Url:         dbPost.Url,
		Description: dbPost.Description.String,
		FeedID:      dbPost.FeedID,
	}
}

func DBPostsToPosts(dbPosts []database.Post) []Post {
	posts := make([]Post, len(dbPosts))
	for i, dbPost := range dbPosts {
		post := DBPostToPost(dbPost)
		posts[i] = post
	}
	return posts
}
