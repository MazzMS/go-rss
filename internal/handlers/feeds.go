package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/MazzMS/go-rss/internal/database"
	"github.com/MazzMS/go-rss/internal/models"
	"github.com/MazzMS/go-rss/internal/utils"
	"github.com/google/uuid"
)

func (cfg *ApiConfig) CreateFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	if cfg.Debug {
		log.Println("CALLING FEED CREATION")
		log.Println()
		defer log.Println("END OF FEED CREATION")
		defer log.Println()
	}

	// types for JSON's input and output
	type parameters struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}

	param := parameters{}

	// decode input
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&param)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprint("failed to decode JSON input:", err))
		return
	}

	// check if some param was null
	if param.Name == "" {
		utils.RespondWithError(w, http.StatusInternalServerError, "name cannot be null")
		return
	}
	if param.URL == "" {
		utils.RespondWithError(w, http.StatusInternalServerError, "url cannot be null")
		return
	}

	currentTime := time.Now().UTC()

	// create feed
	feed, err := cfg.DB.CreateFeed(
		r.Context(),
		database.CreateFeedParams{
			ID:        uuid.New(),
			Name:      param.Name,
			CreatedAt: currentTime,
			UpdatedAt: currentTime,
			Url:       param.URL,
			UserID:    user.ID,
		},
	)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprint("failed to create feed:", err))
		log.Println(err)
		return
	}

	if cfg.Debug {
		log.Printf(
			"successfully recorded feed: %s with URL %v, with UUID %s at %v. For user %v, uuid %v",
			feed.Name, feed.Url, feed.ID, feed.CreatedAt, user.Name, user.ID,
		)
	}

	utils.RespondWithJSON(w, http.StatusCreated, models.DBFeedToFeed(feed))
}

func (cfg *ApiConfig) GetAllFeeds(w http.ResponseWriter, r *http.Request) {
	dbFeeds, err := cfg.DB.GetFeeds(r.Context())
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("failed to get feeds: %v", err))
	}

	utils.RespondWithJSON(w, http.StatusOK, models.DBFeedsToFeeds(dbFeeds))
}
