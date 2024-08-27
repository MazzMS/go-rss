package scraper

import (
	"context"
	"database/sql"
	"encoding/xml"
	"errors"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/MazzMS/go-rss/internal/database"
	"github.com/MazzMS/go-rss/internal/models"
	"github.com/google/uuid"
)

func StartScraping(db *database.Queries, concurrency int, timeBetweenRequest time.Duration) {
	log.Printf("Starting to scrap on %d gorutines every %s duration", concurrency, timeBetweenRequest)

	ticker := time.NewTicker(timeBetweenRequest)

	// this way behave more of a `do... while` meanwhile
	// the `for range ticker.C` has to wait timeBetweenRequest
	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeedsToFetch(
			context.Background(),
			int32(concurrency),
		)
		if err != nil {
			log.Println("could not fetch feeds from db:", err)
			continue
		}

		wg := &sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1)

			go scrapeFeed(db, wg, feed)
		}
		wg.Wait()
	}
}

func scrapeFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
	defer wg.Done()

	_, err := db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		log.Println("failed to mark feed as fetched:", err)
		return
	}

	rssFeed, err := urlToFeed(feed.Url)
	if err != nil {
		log.Println("failed to scrap url:", err)
		return
	}

	totalPosts := len(rssFeed.Channel.Item)
	savedPosts := 0
	newPosts := 0

	for _, RSSItem := range rssFeed.Channel.Item {
		// check if already stored
		exists, err := db.CheckPostByUrl(context.Background(), RSSItem.Link)
		if err != nil {
			log.Println("failed to check if post exists:", err)
			continue
		}
		if exists {
			savedPosts++
			continue
		}

		// parse some values to store them
		publishedAt, err := parseTime(
			RSSItem.PubDate,
			[]string{time.RFC1123Z, time.RFC1123},
		)
		if err != nil {
			log.Printf("failed to parse time %s, err: %v", RSSItem.PubDate, err)
			continue
		}

		currentTime := time.Now().UTC()

		description := sql.NullString{}
		if RSSItem.Description != "" {
			description.String = RSSItem.Description
			description.Valid = true
		}

		_, err = db.CreatePost(
			context.Background(),
			database.CreatePostParams{
				ID:          uuid.New(),
				CreatedAt:   currentTime,
				UpdatedAt:   currentTime,
				Title:       RSSItem.Title,
				Url:         RSSItem.Link,
				Description: description,
				PublishedAt: publishedAt,
				FeedID:      feed.ID,
			},
		)
		if err != nil {
			log.Println("failed to save post:", err)
			continue
		}
		newPosts++
	}
	log.Printf(
		"Feed %s collected, %d posts found, %d new posts and %d were already saved, %d with error/s",
		feed.Name, totalPosts, newPosts, savedPosts, totalPosts-(savedPosts+newPosts),
	)
}

func urlToFeed(url string) (models.RSSFeed, error) {
	httpClient := http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := httpClient.Get(url)
	if err != nil {
		return models.RSSFeed{}, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)

	rssFeed := models.RSSFeed{}

	err = xml.Unmarshal(data, &rssFeed)
	if err != nil {
		return models.RSSFeed{}, err
	}

	return rssFeed, nil
}

func parseTime(possibleTime string, layouts []string) (time.Time, error) {
	for _, layout := range layouts {
		time, err := time.Parse(layout, possibleTime)
		if err == nil {
			return time, nil
		}
	}
	return time.Time{}, errors.New("failed to match with provided layouts")
}
