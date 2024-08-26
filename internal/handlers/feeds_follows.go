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

func (cfg *ApiConfig) CreateFeedFollow(w http.ResponseWriter, r *http.Request, dbUser database.User) {
	if cfg.Debug {
		log.Println("CALLING FEED FOLLOW CREATION")
		log.Println()
		defer log.Println("END OF FEED FOLLOW CREATION")
		defer log.Println()
	}
	// types for JSON's input and output
	type parameters struct {
		FeedID uuid.UUID `json:"feed_id"`
	}

	param := parameters{}

	// decode input
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&param)

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprint("failed to decode JSON input:", err))
		return
	}

	currentTime := time.Now().UTC()
	// create follow
	dbFeedFollow, err := cfg.DB.CreateFeedFollow(
		r.Context(),
		database.CreateFeedFollowParams{
			ID:        uuid.New(),
			CreatedAt: currentTime,
			UpdatedAt: currentTime,
			UserID:    dbUser.ID,
			FeedID:    param.FeedID,
		},
	)

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprint("failed to record follow:", err))
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, models.DBFeedFollowToFeedFollow(dbFeedFollow))
}

func (cfg *ApiConfig) DeleteFeedFollow(w http.ResponseWriter, r *http.Request, dbUser database.User) {
	if cfg.Debug {
		log.Println("CALLING FEED FOLLOW DELETION")
		log.Println()
		defer log.Println("END OF FEED FOLLOW DELETION")
		defer log.Println()
	}
	// get feed_id to unfollow
	pathValue := r.PathValue("feedFollowID")
	if pathValue == "" {
		utils.RespondWithError(w, http.StatusBadRequest, fmt.Sprint("no feed_id provided:", r.URL))
		return
	}

	feedFollowID, err := uuid.Parse(pathValue)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, fmt.Sprint("failed to parse uuid", err))
		return
	}

	// delete follow
	err = cfg.DB.DeleteFeedFollow(
		r.Context(),
		database.DeleteFeedFollowParams{
			UserID: dbUser.ID,
			ID:     feedFollowID,
		},
	)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprint("failed to delete follow:", err))
		return
	}

	if cfg.Debug {
		log.Printf("user %s, uuid %s stopped following a feed", dbUser.Name, dbUser.ID)
	}

	utils.RespondWithJSON(w, http.StatusNoContent, nil)
}

func (cfg *ApiConfig) GetFollowsFeeds(w http.ResponseWriter, r *http.Request, dbUser database.User) {
	if cfg.Debug {
		log.Println("CALLING FEED FOLLOW RECOVER")
		log.Println()
		defer log.Println("END OF FEED FOLLOW RECOVER")
		defer log.Println()
	}
	// get follows
	dbFeedFollows, err := cfg.DB.GetFeedsFollowsByUser(
		r.Context(),
		dbUser.ID,
	)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprint("failed to recover follows:", err))
		return
	}
	if cfg.Debug {
		log.Printf("trying to recover %s's feeds follows:", dbUser.Name)
		log.Println("records from DB")
		for _, dbFeedFollow := range dbFeedFollows {
			log.Printf(" - %v", dbFeedFollow.FeedID)
		}
		log.Println("records parsed")
		for _, feedFollow := range models.DBFeedsFollowsToFeedsFollows(dbFeedFollows) {
			log.Printf(" - %v", feedFollow.FeedID)
		}
	}

	utils.RespondWithJSON(w, http.StatusOK, models.DBFeedsFollowsToFeedsFollows(dbFeedFollows))
}
